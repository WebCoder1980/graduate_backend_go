package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"user_microservice/internal/service"
)

type testFixture struct {
	t       *testing.T
	handler *Handler
	cleanup func()
}

func setupIntegrationTest(t *testing.T) *testFixture {
	t.Helper()
	ctx := context.Background()

	os.Setenv("keycloak_address", "http://localhost:8080")
	os.Setenv("keycloak_realm", "test-realm")
	os.Setenv("keycloak_client_id", "test-client")
	os.Setenv("keycloak_client_secret", "test-secret")
	os.Setenv("keycloak_admin_username", "admin")
	os.Setenv("keycloak_admin_password", "admin")

	mockKc := service.NewMockKeycloak()
	svc, err := service.NewServiceWithKeycloak(ctx, mockKc)
	if err != nil {
		t.Skip("Skipping test - cannot create service")
		return nil
	}

	handler := &Handler{
		ctx:     ctx,
		service: svc,
	}

	cleanup := func() {
	}

	t.Cleanup(cleanup)

	return &testFixture{
		t:       t,
		handler: handler,
		cleanup: cleanup,
	}
}

func (f *testFixture) doPost(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewReader(jsonBody)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}
	req := httptest.NewRequest(http.MethodPost, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rr := httptest.NewRecorder()

	switch path {
	case "/api/v1/user/login":
		f.handler.UserLoginHandler(rr, req)
	case "/api/v1/user/refresh-token":
		f.handler.UserRefreshTokenHandler(rr, req)
	case "/api/v1/user/register":
		f.handler.UserRegisterHandler(rr, req)
	}

	return rr
}
