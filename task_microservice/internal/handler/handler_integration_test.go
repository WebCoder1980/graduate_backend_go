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

func TestIntegration_GetTasks_Success(t *testing.T) {
	f := setupIntegrationTest(t)

	userUuid := "test-user-1"
	token := generateJWTToken(userUuid)

	f.createTask(userUuid, nil, nil, nil, nil)
	f.createTask(userUuid, toIntPtr(100), toIntPtr(200), toStringPtr("jpeg"), toFloat64Ptr(0.8))

	rr := f.doGet("/api/v1/task", token)

	assert.Equal(t, http.StatusOK, rr.Code)

	var tasks []map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &tasks)
	require.NoError(t, err)
	assert.Len(t, tasks, 2)
}

func TestIntegration_GetTasks_Unauthorized(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doGet("/api/v1/task", "")

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "Invalid or missing token")
}

func TestIntegration_GetTasks_InvalidToken(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doGet("/api/v1/task", "invalid-token")

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestIntegration_GetTasks_EmptyList(t *testing.T) {
	f := setupIntegrationTest(t)

	token := generateJWTToken("user-with-no-tasks")

	rr := f.doGet("/api/v1/task", token)

	assert.Equal(t, http.StatusOK, rr.Code)

	var tasks []interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &tasks)
	require.NoError(t, err)
	assert.Len(t, tasks, 0)
}

func TestIntegration_GetTaskById_Success(t *testing.T) {
	f := setupIntegrationTest(t)

	userUuid := "test-user-1"
	token := generateJWTToken(userUuid)

	taskId := f.createTask(userUuid, nil, nil, nil, nil)
	f.createImage(taskId, 1, "test-image", "png", 2)

	rr := f.doGet("/api/v1/task/"+itoa(taskId), token)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(taskId), resp["id"])
}

func TestIntegration_GetTaskById_Unauthorized(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doGet("/api/v1/task/1", "")

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestIntegration_GetTaskById_AccessDenied(t *testing.T) {
	f := setupIntegrationTest(t)

	taskId := f.createTask("owner-user", nil, nil, nil, nil)

	token := generateJWTToken("another-user")

	rr := f.doGet("/api/v1/task/"+itoa(taskId), token)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "Failed to get task images")
}

func TestIntegration_GetTaskById_InvalidId(t *testing.T) {
	f := setupIntegrationTest(t)

	token := generateJWTToken("test-user")

	rr := f.doGet("/api/v1/task/invalid", token)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_GetTaskById_NegativeId(t *testing.T) {
	f := setupIntegrationTest(t)

	token := generateJWTToken("test-user")

	rr := f.doGet("/api/v1/task/-1", token)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestIntegration_PostTask_InvalidContentType(t *testing.T) {
	f := setupIntegrationTest(t)

	token := generateJWTToken("test-user")

	rr := f.doPost("/api/v1/task", token, map[string]string{
		"Content-Type": "application/json",
	})

	assert.Equal(t, http.StatusUnsupportedMediaType, rr.Code)
}

func TestIntegration_PostTask_Unauthorized(t *testing.T) {
	f := setupIntegrationTest(t)

	rr := f.doPost("/api/v1/task", "", map[string]string{
		"Content-Type": "multipart/form-data",
	})

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestIntegration_PostTask_NoFiles(t *testing.T) {
	f := setupIntegrationTest(t)

	token := generateJWTToken("test-user")

	rr := f.doPost("/api/v1/task?width=100&height=200&format=jpeg&quality=0.8", token, map[string]string{
		"Content-Type": "multipart/form-data",
	})

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Contains(t, resp["error"], "Failed to parse form data")
}

func TestIntegration_MethodNotAllowed(t *testing.T) {
	f := setupIntegrationTest(t)

	token := generateJWTToken("test-user")

	req := httptest.NewRequest(http.MethodPut, "/api/v1/task", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	f.handler.TaskHandler(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestIntegration_GetTaskById_CommonStatusInWork(t *testing.T) {
	f := setupIntegrationTest(t)

	userUuid := "test-user-1"
	token := generateJWTToken(userUuid)

	taskId := f.createTask(userUuid, nil, nil, nil, nil)
	f.createImage(taskId, 1, "image1", "png", 1)

	rr := f.doGet("/api/v1/task/"+itoa(taskId), token)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(1), resp["common_status_id"])
}

func TestIntegration_GetTaskById_CommonStatusFailed(t *testing.T) {
	f := setupIntegrationTest(t)

	userUuid := "test-user-1"
	token := generateJWTToken(userUuid)

	taskId := f.createTask(userUuid, nil, nil, nil, nil)
	f.createImage(taskId, 1, "image1", "png", 3)

	rr := f.doGet("/api/v1/task/"+itoa(taskId), token)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(3), resp["common_status_id"])
}

func TestIntegration_GetTaskById_CommonStatusSuccess(t *testing.T) {
	f := setupIntegrationTest(t)

	userUuid := "test-user-1"
	token := generateJWTToken(userUuid)

	taskId := f.createTask(userUuid, nil, nil, nil, nil)
	f.createImage(taskId, 1, "image1", "png", 2)

	rr := f.doGet("/api/v1/task/"+itoa(taskId), token)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, float64(2), resp["common_status_id"])
}

func itoa(i int64) string {
	return strconv.FormatInt(i, 10)
}
