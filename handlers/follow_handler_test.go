package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ChenSongJian/ginstagram/handlers"
	"github.com/ChenSongJian/ginstagram/middlewares"
	"github.com/ChenSongJian/ginstagram/mocks"
	"github.com/ChenSongJian/ginstagram/models"
	"github.com/gin-gonic/gin"
)

func TestFollowUser_MissingToken(t *testing.T) {
	mockFollowService := mocks.NewMockFollowService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.FollowUser(mockFollowService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestFollowUser_InvalidToken(t *testing.T) {
	mockFollowService := mocks.NewMockFollowService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer invalidtoken")

	handlers.FollowUser(mockFollowService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestFollowUser_MissingRequiredField(t *testing.T) {
	mockFollowService := mocks.NewMockFollowService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:       1,
		Username: "test",
		Email:    "test@test.com",
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	jsonBody, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.FollowUser(mockFollowService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "Error:Field validation for 'UserId' failed on the 'required' tag"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestFollowUser_FollowSelf(t *testing.T) {
	mockFollowService := mocks.NewMockFollowService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:       1,
		Username: "test",
		Email:    "test@test.com",
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	jsonBody, err := json.Marshal(map[string]interface{}{
		"user_id": 1,
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.FollowUser(mockFollowService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "can not follow yourself"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestFollowUser_UserNotExists(t *testing.T) {
	mockFollowService := mocks.NewMockFollowService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:       1,
		Username: "test",
		Email:    "test@test.com",
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	jsonBody, err := json.Marshal(map[string]interface{}{
		"user_id": 2,
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.FollowUser(mockFollowService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestFollowUser_AlreadyFollowing(t *testing.T) {
	mockFollowService := mocks.NewMockFollowService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:       1,
		Username: "test",
		Email:    "test@test.com",
	}
	mockFollowService.UserService.Users[testUser.Email] = testUser

	testUser2 := models.User{
		Id:       2,
		Username: "test2",
		Email:    "test2@test.com",
	}
	mockFollowService.UserService.Users[testUser2.Email] = testUser2
	mockFollowService.Follow[testUser.Id] = append(mockFollowService.Follow[testUser.Id], testUser2.Id)
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	jsonBody, err := json.Marshal(map[string]interface{}{
		"user_id": testUser2.Id,
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.FollowUser(mockFollowService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "already following"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestFollowUser_Success(t *testing.T) {
	mockFollowService := mocks.NewMockFollowService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:       1,
		Username: "test",
		Email:    "test@test.com",
	}
	mockFollowService.UserService.Users[testUser.Email] = testUser

	testUser2 := models.User{
		Id:       2,
		Username: "test2",
		Email:    "test2@test.com",
	}
	mockFollowService.UserService.Users[testUser2.Email] = testUser2
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	jsonBody, err := json.Marshal(map[string]interface{}{
		"user_id": testUser2.Id,
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.FollowUser(mockFollowService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	if _, ok := mockFollowService.Follow[testUser.Id]; !ok {
		t.Errorf("Expected follow map to contain key %d, got nil", testUser.Id)
	}
	found := false
	for _, followeeId := range mockFollowService.Follow[testUser.Id] {
		if followeeId == testUser2.Id {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected follow map to contain value %d, got nil", testUser2.Id)
	}
}
