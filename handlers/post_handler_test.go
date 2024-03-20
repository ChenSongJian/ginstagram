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

func TestCreatePost_MissingToken(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mocks.NewMockMediaService()

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.CreatePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}

}

func TestCreatePost_InvalidToken(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mocks.NewMockMediaService()

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer invalidtoken")

	handlers.CreatePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}

}

func TestCreatePost_MissingRequiredField(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mocks.NewMockMediaService()

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
		"title": "Test Post",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}

	context.Request, _ = http.NewRequest("GET", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreatePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "Error:Field validation for 'Content' failed on the 'required' tag"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestCreatePost_NoMedia(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mocks.NewMockMediaService()

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
		"title":   "Test Post",
		"content": "Test Content",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}

	context.Request, _ = http.NewRequest("GET", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreatePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "please upload at least one and no more than 9 media"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestCreatePost_TooManyMedia(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mocks.NewMockMediaService()

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
		"title":   "Test Post",
		"content": "Test Content",
		"media":   []string{"m0", "m1", "m2", "m3", "m4", "m5", "m6", "m7", "m8", "m9"},
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}

	context.Request, _ = http.NewRequest("GET", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreatePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "please upload at least one and no more than 9 media"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestCreatePost_Success(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mocks.NewMockMediaService()

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
		"title":   "Test Post",
		"content": "Test Content",
		"media":   []string{"m0", "m1"},
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}

	context.Request, _ = http.NewRequest("GET", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreatePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "post created successfully"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}
