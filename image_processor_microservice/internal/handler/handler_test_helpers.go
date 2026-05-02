package handler

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"graduate_backend_image_processor_microservice/internal/kafkaproducer"
	"graduate_backend_image_processor_microservice/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testRSAPrivateKey *rsa.PrivateKey
	testRSAPublicKey  *rsa.PublicKey
)

func init() {
	var err error
	testRSAPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	testRSAPublicKey = &testRSAPrivateKey.PublicKey
}

type testFixture struct {
	t           *testing.T
	handler     *Handler
	db          *sql.DB
	minioClient *minio.Client
	containers  *testContainers
	cleanup     func()
}

type testContainers struct {
	postgresContainer testcontainers.Container
	minioContainer    testcontainers.Container
	postgresURL       string
	minioEndpoint     string
}

func setupIntegrationTest(t *testing.T) *testFixture {
	t.Helper()
	ctx := context.Background()

	tc := &testContainers{}

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)
	tc.postgresContainer = pgContainer

	postgresURL, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)
	tc.postgresURL = postgresURL

	pgHost, err := pgContainer.Host(ctx)
	require.NoError(t, err)
	pgPort, err := pgContainer.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	os.Setenv("postgresql_host", pgHost)
	os.Setenv("postgresql_port", pgPort.Port())
	os.Setenv("postgresql_user", "testuser")
	os.Setenv("postgresql_password", "testpass")
	os.Setenv("postgresql_dbname", "testdb")

	minioContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "minio/minio:latest",
			ExposedPorts: []string{"9000/tcp"},
			Env: map[string]string{
				"MINIO_ROOT_USER":     "minioadmin",
				"MINIO_ROOT_PASSWORD": "minioadmin",
			},
			Cmd:        []string{"server", "/data"},
			WaitingFor: wait.ForHTTP("/minio/health/live").WithPort("9000/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)
	tc.minioContainer = minioContainer

	minioHost, err := minioContainer.Host(ctx)
	require.NoError(t, err)
	minioPort, err := minioContainer.MappedPort(ctx, "9000/tcp")
	require.NoError(t, err)
	tc.minioEndpoint = net.JoinHostPort(minioHost, minioPort.Port())

	os.Setenv("minio_address", tc.minioEndpoint)
	os.Setenv("minio_access_key_id", "minioadmin")
	os.Setenv("minio_secret_access_key", "minioadmin")
	os.Setenv("minio_token", "")
	os.Setenv("minio_use_ssl", "false")

	pubKeyBytes := rsaPublicKeyToPEM(testRSAPublicKey)
	os.Setenv("keycloak_publickey", pubKeyBytes)

	mockProducer := kafkaproducer.NewMockProducer()
	svc, err := service.NewServiceWithProducer(ctx, mockProducer)
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		_ = minioContainer.Terminate(ctx)
		t.Skip("Skipping test - cannot create service")
		return nil
	}

	handler := NewHandlerWithService(svc)

	db, err := sql.Open("pgx", postgresURL)
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS image_status(
			id BIGINT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
		INSERT INTO image_status(id, name) VALUES (1, 'Обрабатывается'), (2, 'Успех'), (3, 'Ошибка')
		ON CONFLICT DO NOTHING;

		CREATE TABLE IF NOT EXISTS image(
			id BIGSERIAL PRIMARY KEY,
			task_id BIGINT NOT NULL,
			position INT NOT NULL,
			name TEXT NOT NULL,
			format TEXT NOT NULL,
			status_id BIGINT REFERENCES image_status(id) NOT NULL,
			end_dt TIMESTAMP NULL,
			user_uuid TEXT NOT NULL,
			CONSTRAINT uq_image_task_id_position UNIQUE (task_id, position)
		);
	`)
	require.NoError(t, err)

	minioClient, err := minio.New(tc.minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	require.NoError(t, err)

	err = minioClient.MakeBucket(ctx, "image-result", minio.MakeBucketOptions{})
	if err != nil {
		exists, _ := minioClient.BucketExists(ctx, "image-result")
		require.True(t, exists, "bucket should exist or be created")
	}

	err = minioClient.MakeBucket(ctx, "image-source", minio.MakeBucketOptions{})
	if err != nil {
		exists, _ := minioClient.BucketExists(ctx, "image-source")
		require.True(t, exists, "bucket should exist or be created")
	}

	cleanup := func() {
		_, _ = db.ExecContext(ctx, "TRUNCATE TABLE image, image_status RESTART IDENTITY CASCADE")
		_ = db.Close()
		_ = pgContainer.Terminate(ctx)
		_ = minioContainer.Terminate(ctx)
	}

	t.Cleanup(cleanup)

	return &testFixture{
		t:           t,
		handler:     handler,
		db:          db,
		minioClient: minioClient,
		containers:  tc,
		cleanup:     cleanup,
	}
}

func generateJWTToken(sub string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(testRSAPrivateKey)
	return tokenString
}

func rsaPublicKeyToPEM(pub *rsa.PublicKey) string {
	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return ""
	}
	pemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	return string(pem.EncodeToMemory(pemBlock))
}

func (f *testFixture) createImage(taskId int64, position int, name, format string, statusId int64, userUuid string) int64 {
	var id int64
	err := f.db.QueryRowContext(context.Background(),
		"INSERT INTO image (task_id, position, name, format, status_id, user_uuid) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		taskId, position, name, format, statusId, userUuid,
	).Scan(&id)
	require.NoError(f.t, err)

	minioFilename := fmt.Sprintf("%d_%d.%s", taskId, position, format)
	_, err = f.minioClient.PutObject(context.Background(), "image-result", minioFilename,
		strings.NewReader("test-image-data"), int64(len("test-image-data")), minio.PutObjectOptions{})
	require.NoError(f.t, err)

	return id
}

func (f *testFixture) doGet(path, token string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()

	id := path[len("/api/v1/image-processor/"):]
	req.SetPathValue("id", id)
	f.handler.ImageIdGet(rr, req)

	return rr
}

func getAvailablePort(t *testing.T) string {
	l, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	port := strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
	_ = l.Close()
	return port
}
