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

func TestCreateComment_MissingToken(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("POST", "/", nil)

	handlers.CreateComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestCreateComment_InvalidPostId(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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
			Key:   "postId",
			Value: "invalid_id",
		},
	}

	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreateComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid post id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestCreateComment_PostNotFound(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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
			Key:   "postId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreateComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "post not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestCreateComment_NoPermission(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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

	mockCommentService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 2,
	}
	mockCommentService.UserService.Users["test@email.com"] = models.User{
		Id:        2,
		IsPrivate: true,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreateComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "post is private and you are not following the author"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestCreateComment_MissingRequiredField(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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

	mockCommentService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 1,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}
	jsonBody, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}

	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreateComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "Error:Field validation for 'Content' failed on the 'required' tag"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestCreateComment_SuccessOwnPost(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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

	mockCommentService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 1,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}
	jsonBody, err := json.Marshal(map[string]interface{}{
		"content": "test comment",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}

	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreateComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "comment created successfully"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}

	commentRecord := mocks.CommentRecord{
		UserId:  1,
		PostId:  1,
		Content: "test comment",
	}
	found := false
	for _, comment := range mockCommentService.Comments {
		if comment == commentRecord {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Comment not added")
	}
}

func TestCreateComment_SuccessPublicPost(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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

	mockCommentService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 2,
	}
	mockCommentService.UserService.Users["test@email.com"] = models.User{
		Id:        2,
		IsPrivate: false,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}
	jsonBody, err := json.Marshal(map[string]interface{}{
		"content": "test comment",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}

	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreateComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "comment created successfully"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}

	commentRecord := mocks.CommentRecord{
		UserId:  1,
		PostId:  1,
		Content: "test comment",
	}
	found := false
	for _, comment := range mockCommentService.Comments {
		if comment == commentRecord {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Comment not added")
	}
}

func TestCreateComment_SuccessFollowingPost(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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

	mockCommentService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 2,
	}
	mockCommentService.UserService.Users["test@email.com"] = models.User{
		Id:        2,
		IsPrivate: true,
	}
	mockCommentService.FollowService.Follows[1] = mocks.FollowRecord{
		FolloweeId: 2,
		FollowerId: 1,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}
	jsonBody, err := json.Marshal(map[string]interface{}{
		"content": "test comment",
	})
	if err != nil {
		t.Errorf("Error marshaling request body: %v", err)
		return
	}

	context.Request, _ = http.NewRequest("POST", "/", bytes.NewReader(jsonBody))
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.CreateComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "comment created successfully"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}

	commentRecord := mocks.CommentRecord{
		UserId:  1,
		PostId:  1,
		Content: "test comment",
	}
	found := false
	for _, comment := range mockCommentService.Comments {
		if comment == commentRecord {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Comment not added")
	}
}

func TestDeleteComment_MissingToken(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("DELETE", "/", nil)

	handlers.DeleteComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeleteComment_InvalidPostId(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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
			Key:   "postId",
			Value: "invalid_id",
		},
		{
			Key:   "commentId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.DeleteComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid post id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeleteComment_InvalidCommentId(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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
			Key:   "postId",
			Value: "1",
		},
		{
			Key:   "commentId",
			Value: "invalid_id",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.DeleteComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid comment id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeleteComment_CommentNotFound(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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
			Key:   "postId",
			Value: "1",
		},
		{
			Key:   "commentId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.DeleteComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "comment not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeleteComment_NoPermission(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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

	mockCommentService.Comments[1] = mocks.CommentRecord{
		UserId: 2,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
		{
			Key:   "commentId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.DeleteComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "you are not the author of the comment"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestDeleteComment_Success(t *testing.T) {
	mockCommentService := mocks.NewMockCommentService()
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
	commentId := 1
	mockCommentService.Comments[commentId] = mocks.CommentRecord{
		UserId: 1,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
		{
			Key:   "commentId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.DeleteComment(mockCommentService.UserService, mockCommentService.FollowService,
		mockCommentService.PostService, mockCommentService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	if _, ok := mockCommentService.Comments[commentId]; ok {
		t.Errorf("Comment not deleted")
	}
}
