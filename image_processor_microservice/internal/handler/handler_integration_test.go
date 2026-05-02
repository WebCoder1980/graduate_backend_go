package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_ImageIdGet_Success(t *testing.T) {
	f := setupIntegrationTest(t)

	userUuid := "test-user-1"
	token := generateJWTToken(userUuid)

	imageId := f.createImage(1, 1, "test-image", "png", 2, userUuid)

	rr := f.doGet("/api/v1/image-processor/"+itoa(imageId), token)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestIntegration_ImageIdGet_Unauthorized(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doGet("/api/v1/image-processor/1", "")

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "Invalid or missing token")
}

func TestIntegration_ImageIdGet_InvalidToken(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doGet("/api/v1/image-processor/1", "invalid-token")

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_ImageIdGet_InvalidId(t *testing.T) {
	f := setupIntegrationTest(t)

	token := generateJWTToken("test-user")

	rr := f.doGet("/api/v1/image-processor/invalid", token)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_ImageIdGet_NegativeId(t *testing.T) {
	f := setupIntegrationTest(t)

	token := generateJWTToken("test-user")

	rr := f.doGet("/api/v1/image-processor/-1", token)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_ImageIdGet_AccessDenied(t *testing.T) {
	f := setupIntegrationTest(t)

	imageId := f.createImage(1, 1, "test-image", "png", 2, "owner-user")

	token := generateJWTToken("another-user")

	rr := f.doGet("/api/v1/image-processor/"+itoa(imageId), token)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "Failed to get image")
}

func TestIntegration_ImageIdGet_NotFound(t *testing.T) {
	f := setupIntegrationTest(t)

	token := generateJWTToken("test-user")

	rr := f.doGet("/api/v1/image-processor/999", token)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestIntegration_MethodNotAllowed(t *testing.T) {
	f := setupIntegrationTest(t)

	token := generateJWTToken("test-user")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/image-processor/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	f.handler.ImageIdHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func itoa(i int64) string {
	return strconv.FormatInt(i, 10)
}
