package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestUpdateUser_MissingToken(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()

	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "1",
		},
	}
	jsonBody, err := json.Marshal(map[string]string{
		"username": "username",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("PUT", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")
	middlewares.AuthMiddleware()(context)
	handlers.UpdateUser(mockUserService)(context)
	if response.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUpdateUser_InvalidUserId(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
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

	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "user",
		},
	}
	jsonBody, err := json.Marshal(map[string]string{
		"username": "username",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("PUT", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.UpdateUser(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid user id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUpdateUser_DifferentUserIdAndTokenUserId(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
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

	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "2",
		},
	}
	jsonBody, err := json.Marshal(map[string]string{
		"username": "username",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("PUT", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.UpdateUser(mockUserService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "no permission to update user info"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUpdateUser_MiisingRequiredField(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
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

	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "1",
		},
	}
	jsonBody, err := json.Marshal(map[string]string{
		"username": "username",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("PUT", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.UpdateUser(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "Error:Field validation for 'Bio' failed on the 'required' tag"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	testUser := models.User{
		Id:       1,
		Username: "test",
		Email:    "test@test.com",
	}
	mockUserService.Users["test@email.com"] = models.User{
		Id: 2,
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}

	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "1",
		},
	}
	jsonBody, err := json.Marshal(map[string]interface{}{
		"username":          "username",
		"bio":               "bio",
		"profile_image_url": "profile_image_url",
		"is_private":        true,
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("PUT", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.UpdateUser(mockUserService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "user not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUpdateUser_Success(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	testUser := models.User{
		Id:       1,
		Username: "test",
		Email:    "test@test.com",
	}
	mockUserService.Users[testUser.Email] = testUser
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}

	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "1",
		},
	}
	jsonBody, err := json.Marshal(map[string]interface{}{
		"username":          "username",
		"bio":               "bio",
		"profile_image_url": "profile_image_url",
		"is_private":        true,
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("PUT", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.UpdateUser(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	var responseBody handlers.UserResponse
	err = json.Unmarshal(response.Body.Bytes(), &responseBody)
	if err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
		return
	}
	if responseBody.Id != testUser.Id || responseBody.Username != "username" ||
		responseBody.Email != testUser.Email || responseBody.Bio != "bio" ||
		responseBody.ProfileImageUrl != "profile_image_url" || responseBody.IsPrivate != true {
		t.Errorf("Expected response body %v, got %v", testUser, responseBody)
	}
}

func TestDeleteUser_MissingToken(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "1",
		},
	}
	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Content-Type", "application/json")
	middlewares.AuthMiddleware()(context)
	handlers.DeleteUser(mockUserService)(context)
	if response.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeleteUser_InvalidUserId(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
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

	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "user",
		},
	}
	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Content-Type", "application/json")
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.DeleteUser(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid user id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeleteUser_DifferentUserIdAndTokenUserId(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
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

	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "2",
		},
	}
	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Content-Type", "application/json")
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.DeleteUser(mockUserService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "no permission to delete user info"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeleteUser_UserNotFound(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
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

	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "1",
		},
	}
	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Content-Type", "application/json")
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.DeleteUser(mockUserService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "user not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeleteUser_Success(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:       1,
		Username: "test",
		Email:    "test@test.com",
	}
	mockUserService.Users[testUser.Email] = testUser
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}

	context.Params = []gin.Param{
		{
			Key:   "userId",
			Value: "1",
		},
	}
	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Content-Type", "application/json")
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.DeleteUser(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	if _, ok := mockUserService.Users[testUser.Email]; ok {
		t.Errorf("User %v not deleted", testUser)
	}
}

func TestGetCurrentUserInfo_MissingToken(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.GetCurrentUserInfo(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestGetCurrentUserInfo_UserNotFound(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
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
	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.GetCurrentUserInfo(mockUserService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "user not found\""
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestGetCurrentUserInfo_Success(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:       1,
		Username: "test",
		Email:    "test@test.com",
	}
	mockUserService.Users[testUser.Email] = testUser
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.GetCurrentUserInfo(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	var userResponse handlers.UserResponse
	err = json.Unmarshal(response.Body.Bytes(), &userResponse)
	if err != nil {
		t.Errorf("Error unmarshaling response body: %v", err)
		return
	}
	if userResponse.Id != testUser.Id {
		t.Errorf("Expected user id %d, got %d", testUser.Id, userResponse.Id)
	}
	if userResponse.Username != testUser.Username {
		t.Errorf("Expected user username %s, got %s", testUser.Username, userResponse.Username)
	}
	if userResponse.Email != testUser.Email {
		t.Errorf("Expected user email %s, got %s", testUser.Email, userResponse.Email)
	}
}

func TestLoginUser_MissingRequiredField(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	jsonBody, err := json.Marshal(map[string]string{
		"email": "example@email.com",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")
	handlers.LoginUser(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "Error:Field validation for 'Password' failed on the 'required' tag"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLoginUser_InvalidEmail(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	jsonBody, err := json.Marshal(map[string]string{
		"email":    "email.com",
		"password": "Password123",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")
	handlers.LoginUser(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "Error:Field validation for 'Email' failed on the 'email' tag"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLoginUser_EmailNotFound(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	jsonBody, err := json.Marshal(map[string]string{
		"email":    "email@example.com",
		"password": "Password123",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")
	handlers.LoginUser(mockUserService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "user not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	jsonBody, err := json.Marshal(map[string]string{
		"username": "username",
		"email":    "email@example.com",
		"password": "Password123",
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

	jsonBody, err = json.Marshal(map[string]string{
		"email":    "email@example.com",
		"password": "Password456",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	response = httptest.NewRecorder()
	context, _ = gin.CreateTestContext(response)
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")
	handlers.LoginUser(mockUserService)(context)
	if response.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, response.Code)
	}
	expectedResponseBodyString := "invalid password"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLoginUser_Success(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	jsonBody, err := json.Marshal(map[string]string{
		"username": "username",
		"email":    "email@example.com",
		"password": "Password123",
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

	jsonBody, err = json.Marshal(map[string]string{
		"email":    "email@example.com",
		"password": "Password123",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}
	response = httptest.NewRecorder()
	context, _ = gin.CreateTestContext(response)
	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Content-Type", "application/json")
	handlers.LoginUser(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "\"token\":\"ey"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLogoutUser_MissingToken(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("POST", "/", nil)

	handlers.LogoutUser(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLogoutUser_TokenUserNotFound(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
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
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	middlewares.AuthMiddleware()(context)

	handlers.LogoutUser(mockUserService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "user not found\""
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLogoutUser_Success(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
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
	mockUserService.Users[testUser.Email] = testUser
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	middlewares.AuthMiddleware()(context)

	handlers.LogoutUser(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "\"token\":\"ey"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestRefreshToken_MissingToken(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.RefreshToken(mockUserService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestResfreshToken_UserNotFound(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
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
	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.RefreshToken(mockUserService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "user not found\""
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestResfreshToken_Success(t *testing.T) {
	mockUserService := mocks.NewMockUserService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:       1,
		Username: "test",
		Email:    "test@test.com",
	}
	mockUserService.Users[testUser.Email] = testUser
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)
	middlewares.AuthMiddleware()(context)
	handlers.RefreshToken(mockUserService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "\"token\":\"ey"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}
