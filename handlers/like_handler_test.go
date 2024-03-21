package handlers_test

import (
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

func TestLikePost_MissingToken(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("POST", "/", nil)

	handlers.LikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikePost_InvalidPostId(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

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
	handlers.LikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid post id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikePost_PostNotFound(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

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
	handlers.LikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "post not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikePost_UserNotFound(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

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
	mockLikeService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 1,
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
	handlers.LikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "user not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikePost_NoPermission(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

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
	mockLikeService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 2,
	}
	mockLikeService.UserService.Users["test@email.com"] = models.User{
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
	handlers.LikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "post is private and you are not following the author"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikePost_SuccessOwnPost(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

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
	mockLikeService.UserService.Users["test@test.com"] = testUser
	mockLikeService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 1,
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
	handlers.LikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	for _, like := range mockLikeService.PostLikes {
		if like.UserId == 1 && like.PostId == 1 {
			return
		}
	}
	t.Errorf("Like not found in post likes")
}

func TestLikePost_SuccessPublicPost(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

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
	mockLikeService.UserService.Users["test@test.com"] = testUser
	mockLikeService.UserService.Users["test2@test.com"] = models.User{
		Id:        2,
		IsPrivate: false,
	}
	mockLikeService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 2,
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
	handlers.LikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	for _, like := range mockLikeService.PostLikes {
		if like.UserId == 1 && like.PostId == 1 {
			return
		}
	}
	t.Errorf("Like not found in post likes")
}

func TestLikePost_SuccessPrivateFollowingPost(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

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
	mockLikeService.UserService.Users["test@test.com"] = testUser
	mockLikeService.UserService.Users["test2@test.com"] = models.User{
		Id:        2,
		IsPrivate: true,
	}
	mockLikeService.FollowService.Follows[1] = mocks.FollowRecord{
		FollowerId: 1,
		FolloweeId: 2,
	}
	mockLikeService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 2,
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
	handlers.LikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	for _, like := range mockLikeService.PostLikes {
		if like.UserId == 1 && like.PostId == 1 {
			return
		}
	}
	t.Errorf("Like not found in post likes")
}

func TestLikePost_DuplicatedLike(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

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
	mockLikeService.UserService.Users["test@test.com"] = testUser
	mockLikeService.UserService.Users["test2@test.com"] = models.User{
		Id:        2,
		IsPrivate: true,
	}
	mockLikeService.FollowService.Follows[1] = mocks.FollowRecord{
		FollowerId: 1,
		FolloweeId: 2,
	}
	mockLikeService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 2,
	}
	mockLikeService.PostLikes[1] = mocks.PostLikeRecord{
		UserId: 1,
		PostId: 1,
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
	handlers.LikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "already liked"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}
