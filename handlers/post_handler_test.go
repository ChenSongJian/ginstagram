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
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("POST", "/", nil)

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
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("POST", "/", nil)
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
	mockMediaService := mockPostService.MediaService

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

	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
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
	mockMediaService := mockPostService.MediaService

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

	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
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
	mockMediaService := mockPostService.MediaService

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

	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
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
	mockMediaService := mockPostService.MediaService

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
	testMediaUrls := []string{"m0", "m1"}
	jsonBody, err := json.Marshal(map[string]interface{}{
		"title":   "Test Post",
		"content": "Test Content",
		"media":   testMediaUrls,
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}

	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreatePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}

	postRecord := mocks.PostRecord{
		Title:   "Test Post",
		Content: "Test Content",
		UserId:  1,
	}
	found := false
	for _, post := range mockPostService.Posts {
		if postRecord == post {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected post %v, not found", postRecord)
	}

	mediaUrls := []string{}
	for _, media := range mockMediaService.Media {
		mediaUrls = append(mediaUrls, media.Url)
	}
	if len(mediaUrls) != len(testMediaUrls) {
		t.Errorf("Expected %d media, got %d", len(testMediaUrls), len(mediaUrls))
	}
	for _, mediaUrl := range mediaUrls {
		found = false
		for _, testMediaUrl := range testMediaUrls {
			if mediaUrl == testMediaUrl {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected media %s, not found", mediaUrl)
		}
	}
}

func TestDeletePost_MissingToken(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("DELETE", "/", nil)

	handlers.CreatePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}

}

func TestDeletePost_InvalidToken(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
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

func TestDeletePost_InvalidPostId(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "postId",
		},
	}
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

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.DeletePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid post id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeletePost_PostNotFound(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}
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

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.DeletePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "post not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeletePost_DifferentPostUserIdAndTokenUserId(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}
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
	mockPostService.Posts[1] = mocks.PostRecord{
		Title:   "Test Post",
		Content: "Test Content",
		UserId:  2,
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.DeletePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "no permission to delete this post"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeletePost_Success(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}
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
	mockPostService.Posts[1] = mocks.PostRecord{
		Title:   "Test Post",
		Content: "Test Content",
		UserId:  1,
	}
	mockMediaService.Media[1] = mocks.MediaRecord{
		Url:    "m0",
		PostId: 1,
	}
	mockMediaService.Media[2] = mocks.MediaRecord{
		Url:    "m1",
		PostId: 2,
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.DeletePost(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}

	if _, ok := mockPostService.Posts[1]; ok {
		t.Errorf("Expected post %d to be deleted, but it is still present", 1)
	}

	for _, media := range mockMediaService.Media {
		if media.PostId == 1 {
			t.Errorf("Expected media %d to be deleted, but it is still present", 1)
		}
	}

	if len(mockMediaService.Media) != 1 {
		t.Errorf("Expected %d media, got %d", 1, len(mockMediaService.Media))
	}
}
