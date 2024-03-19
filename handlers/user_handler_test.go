package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ChenSongJian/ginstagram/handlers"
	"github.com/ChenSongJian/ginstagram/mocks"
	"github.com/ChenSongJian/ginstagram/models"
	"github.com/gin-gonic/gin"
)

func TestRegisterUser_MissingRequiredField(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	jsonBody, err := json.Marshal(map[string]string{
		"username": "username",
		"password": "Password123",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")

	handlers.RegisterUser(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "Error:Field validation for 'Email' failed on the 'required' tag"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected error message %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestRegisterUser_InvalidEmailFormat(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	jsonBody, err := json.Marshal(map[string]string{
		"username": "username",
		"password": "Password123",
		"Email":    "Email",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")

	handlers.RegisterUser(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "Error:Field validation for 'Email' failed on the 'email' tag"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected error message %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestRegisterUser_PasswordNotComplex(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	jsonBody, err := json.Marshal(map[string]string{
		"username": "username",
		"password": "password123",
		"Email":    "Email@example.com",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")

	handlers.RegisterUser(mockUserService)(context)
	handlers.RegisterUser(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "password must be complex"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected error message %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestRegisterUser_Success(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	jsonBody, err := json.Marshal(map[string]string{
		"username": "Username",
		"password": "Password123",
		"Email":    "Email@example.com",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")

	handlers.RegisterUser(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "User registered successfully"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestRegisterUser_DuplicatedEmail(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	duplicateUser := models.User{
		Username:     "Username",
		PasswordHash: "Password123",
		Email:        "Email@example.com",
	}
	mockUserService.Create(duplicateUser)

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	jsonBody, err := json.Marshal(map[string]string{
		"username": "Username",
		"password": "Password123",
		"email":    "Email@example.com",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")

	handlers.RegisterUser(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "email already exists"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListUsers_NoUser(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	handlers.ListUsers(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "\"total_records\":0"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListUsers_Success(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	dummyUser := models.User{
		Id:           1,
		Username:     "XXXXXXXX",
		PasswordHash: "PasswordHash",
		Email:        "Email@example.com",
	}
	mockUserService.Create(dummyUser)
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	handlers.ListUsers(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "\"total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListUsers_SuccessWithPagination(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	dummyUser := models.User{
		Id:           1,
		Username:     "XXXXXXXX",
		PasswordHash: "PasswordHash",
		Email:        "Email@example.com",
	}
	mockUserService.Create(dummyUser)
	dummyUser.Id = 2
	dummyUser.Email = "Email2@example.com"
	mockUserService.Create(dummyUser)
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request, _ = http.NewRequest("GET", "/?pageNum=1&pageSize=1", nil)

	handlers.ListUsers(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "\"total_records\":2"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListUsers_WithKeyword(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	dummyUser := models.User{
		Id:           1,
		Username:     "user1",
		PasswordHash: "PasswordHash",
		Email:        "Email@example.com",
	}
	mockUserService.Create(dummyUser)
	dummyUser.Id = 2
	dummyUser.Username = "user2"
	dummyUser.Email = "Email2@example.com"
	mockUserService.Create(dummyUser)
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request, _ = http.NewRequest("GET", "/?keyword=user2", nil)

	handlers.ListUsers(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "\"total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestGetUserById_InvalidUserId(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "invalid_id",
		},
	}

	handlers.GetUserById(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid user id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestGetUserById_UserNotFound(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	dummyUser := models.User{
		Id:           1,
		Username:     "Username",
		PasswordHash: "PasswordHash",
		Email:        "Email@example.com",
	}
	mockUserService.Create(dummyUser)
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "2",
		},
	}

	handlers.GetUserById(mockUserService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "user not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestGetUserById_Success(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	dummyUser := models.User{
		Id:           1,
		Username:     "Username",
		PasswordHash: "PasswordHash",
		Email:        "Email@example.com",
	}
	mockUserService.Create(dummyUser)
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "1",
		},
	}

	handlers.GetUserById(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	var responseBody handlers.UserResponse
	err := json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
		return
	}
	if responseBody.Id != dummyUser.Id ||
		responseBody.Username != dummyUser.Username ||
		responseBody.Email != dummyUser.Email {
		t.Errorf("Expected response body %v, got %v", dummyUser, responseBody)
	}
}
