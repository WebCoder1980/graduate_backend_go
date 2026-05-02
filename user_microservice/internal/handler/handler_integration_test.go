package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_UserLogin_InvalidContentType(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/login", map[string]string{
		"username": "test",
		"password": "test",
	}, map[string]string{
		"Content-Type": "text/plain",
	})

	assert.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
}

func TestIntegration_UserLogin_NoBody(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/login", nil, map[string]string{
		"Content-Type": "application/json",
	})

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_UserLogin_MethodNotAllowed(t *testing.T) {
	f := setupIntegrationTest(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user/login", nil)
	rr := httptest.NewRecorder()
	f.handler.UserLoginHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestIntegration_UserRegister_InvalidContentType(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/register", map[string]string{
		"username": "test",
		"password": "test123",
	}, map[string]string{
		"Content-Type": "text/plain",
	})

	assert.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
}

func TestIntegration_UserRegister_MissingUsername(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/register", map[string]string{
		"password": "test123",
	}, map[string]string{
		"Content-Type": "application/json",
	})

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "Username and password are required")
}

func TestIntegration_UserRegister_MissingPassword(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/register", map[string]string{
		"username": "testuser",
	}, map[string]string{
		"Content-Type": "application/json",
	})

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "Username and password are required")
}

func TestIntegration_UserRegister_ShortPassword(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/register", map[string]string{
		"username": "testuser",
		"password": "123",
	}, map[string]string{
		"Content-Type": "application/json",
	})

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "Password must be at least 6 characters")
}

func TestIntegration_UserRegister_ShortUsername(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/register", map[string]string{
		"username": "ab",
		"password": "test123",
	}, map[string]string{
		"Content-Type": "application/json",
	})

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "Username must be at least 3 characters")
}

func TestIntegration_RefreshToken_InvalidContentType(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/refresh-token", map[string]string{
		"refresh_token": "test-token",
	}, map[string]string{
		"Content-Type": "text/plain",
	})

	assert.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
}

func TestIntegration_RefreshToken_MissingToken(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/refresh-token", map[string]string{}, map[string]string{
		"Content-Type": "application/json",
	})

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "Refresh token is required")
}

func TestIntegration_UserLogin_Success(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/login", map[string]string{
		"username": "test",
		"password": "test",
	}, map[string]string{})

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["access_token"])
}

func TestIntegration_UserRegister_Success(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/register", map[string]string{
		"username": "testuser",
		"password": "test123",
	}, map[string]string{})

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestIntegration_RefreshToken_Success(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/user/refresh-token", map[string]string{
		"refresh_token": "mock-refresh-token",
	}, map[string]string{})

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.NotEmpty(t, resp["access_token"])
}
