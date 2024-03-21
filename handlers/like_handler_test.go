package handlers_test

import (
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

func TestListLikeByPostId_MissingToken(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.ListLikesByPostId(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}
func TestListLikeByPostId_InvalidPostId(t *testing.T) {
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

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListLikesByPostId(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid post id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListLikeByPostId_PostNotFound(t *testing.T) {
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

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListLikesByPostId(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "post not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListLikeByPostId_NoPermission(t *testing.T) {
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
	mockLikeService.UserService.Users["test@test.com"] = models.User{
		Id:        2,
		IsPrivate: true,
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

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListLikesByPostId(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "post is private and you are not following the author"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListLikeByPostId_SuccessOwnPost(t *testing.T) {
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

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListLikesByPostId(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	var responseMap map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &responseMap); err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	if responseMap["likes"] == nil {
		t.Errorf("Expected likes to be present in response")
	}
}

func TestListLikeByPostId_SuccessPublicPost(t *testing.T) {
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
	mockLikeService.UserService.Users["test@test.com"] = models.User{
		Id:        2,
		IsPrivate: false,
	}
	mockLikeService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 1,
	}
	mockLikeService.PostLikes[1] = mocks.PostLikeRecord{
		UserId: 2,
		PostId: 1,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListLikesByPostId(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	var responseMap map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &responseMap); err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	if responseMap["likes"] == nil {
		t.Errorf("Expected likes to be present in response")
	}
}

func TestListLikeByPostId_SuccessPrivateFollowingPost(t *testing.T) {
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
	mockLikeService.UserService.Users["test@test.com"] = models.User{
		Id:        2,
		IsPrivate: true,
	}
	mockLikeService.FollowService.Follows[1] = mocks.FollowRecord{
		FollowerId: 1,
		FolloweeId: 2,
	}
	mockLikeService.PostService.Posts[1] = mocks.PostRecord{
		UserId: 1,
	}
	mockLikeService.PostLikes[1] = mocks.PostLikeRecord{
		UserId: 2,
		PostId: 1,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListLikesByPostId(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	var responseMap map[string]interface{}
	if err := json.Unmarshal(response.Body.Bytes(), &responseMap); err != nil {
		t.Errorf("Error unmarshalling response body: %v", err)
	}
	if responseMap["likes"] == nil {
		t.Errorf("Expected likes to be present in response")
	}
}

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

func TestUnlikePost_MissingToken(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("DELETE", "/", nil)

	handlers.UnlikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUnlikePost_InvalidPostLikeId(t *testing.T) {
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
			Key:   "postLikeId",
			Value: "invalid_id",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.UnlikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid post like id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUnlikePost_PostLikeNotFound(t *testing.T) {
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
			Key:   "postLikeId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.UnlikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "post like not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUnlikePost_NoPermission(t *testing.T) {
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
	mockLikeService.PostLikes[1] = mocks.PostLikeRecord{
		UserId: 2,
		PostId: 1,
	}

	context.Params = []gin.Param{
		{
			Key:   "postLikeId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.UnlikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "no permission to unlike"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUnlikePost_Success(t *testing.T) {
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
	mockLikeService.PostLikes[1] = mocks.PostLikeRecord{
		UserId: 1,
		PostId: 1,
	}

	context.Params = []gin.Param{
		{
			Key:   "postLikeId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.UnlikePost(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	if _, ok := mockLikeService.PostLikes[1]; ok {
		t.Errorf("Post like not deleted")
	}
}

func TestLikeComment_MissingToken(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("POST", "/", nil)

	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikeComment_InvalidPostId(t *testing.T) {
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
	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid post id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikeComment_PostNotFound(t *testing.T) {
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
	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "post not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikeComment_NoPermission(t *testing.T) {
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
	mockLikeService.UserService.Users["test@email.com"] = models.User{
		Id:        2,
		IsPrivate: true,
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
	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "post is private and you are not following the author"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikeComment_InvalidCommentId(t *testing.T) {
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
		{
			Key:   "commentId",
			Value: "invalid",
		},
	}

	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid comment id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikeComment_CommentNotFound(t *testing.T) {
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
		{
			Key:   "commentId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "comment not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikeComment_UserNotFound(t *testing.T) {
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
	mockLikeService.CommentService.Comments[1] = mocks.CommentRecord{
		UserId: 2,
		PostId: 1,
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

	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "user not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikeComment_CommentNotBelongToPost(t *testing.T) {
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
	mockLikeService.CommentService.Comments[1] = mocks.CommentRecord{
		PostId: 2,
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

	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "comment does not belong to the post"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestLikeComment_SuccessOwnPost(t *testing.T) {
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
	mockLikeService.CommentService.Comments[1] = mocks.CommentRecord{
		PostId: 1,
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

	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	if len(mockLikeService.CommentLikes) != 1 {
		t.Errorf("Expected 1 like, got %d", len(mockLikeService.CommentLikes))
	}
}

func TestLikeComment_SuccessPublicPost(t *testing.T) {
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
	mockLikeService.CommentService.Comments[1] = mocks.CommentRecord{
		PostId: 1,
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

	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	if len(mockLikeService.CommentLikes) != 1 {
		t.Errorf("Expected 1 like, got %d", len(mockLikeService.CommentLikes))
	}
}

func TestLikeComment_SuccessPrivateFollowingPost(t *testing.T) {
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
	mockLikeService.CommentService.Comments[1] = mocks.CommentRecord{
		PostId: 1,
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

	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	if len(mockLikeService.CommentLikes) != 1 {
		t.Errorf("Expected 1 like, got %d", len(mockLikeService.CommentLikes))
	}
}

func TestLikeComment_DuplicatedLike(t *testing.T) {
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
	mockLikeService.CommentService.Comments[1] = mocks.CommentRecord{
		PostId: 1,
	}
	mockLikeService.CommentLikes[1] = mocks.CommentLikeRecord{
		UserId:    1,
		CommentId: 1,
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

	context.Request, _ = http.NewRequest("POST", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.LikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.PostService, mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "already like"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUnlikeComment_MissingToken(t *testing.T) {
	mockLikeService := mocks.NewMockLikeService()

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("DELETE", "/", nil)

	handlers.UnlikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUnlikeComment_InvalidPostLikeId(t *testing.T) {
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
			Key:   "commentLikeId",
			Value: "invalid_id",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.UnlikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid post like id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUnlikeComment_PostLikeNotFound(t *testing.T) {
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
			Key:   "commentLikeId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.UnlikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "comment like not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUnlikeComment_NoPermission(t *testing.T) {
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
	mockLikeService.CommentLikes[1] = mocks.CommentLikeRecord{
		UserId:    2,
		CommentId: 1,
	}

	context.Params = []gin.Param{
		{
			Key:   "commentLikeId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.UnlikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "no permission to unlike"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestUnlikeComment_Success(t *testing.T) {
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
	mockLikeService.CommentLikes[1] = mocks.CommentLikeRecord{
		UserId:    1,
		CommentId: 1,
	}

	context.Params = []gin.Param{
		{
			Key:   "commentLikeId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("DELETE", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.UnlikeComment(mockLikeService.UserService, mockLikeService.FollowService,
		mockLikeService.CommentService, mockLikeService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	if _, ok := mockLikeService.CommentLikes[1]; ok {
		t.Errorf("Post like not deleted")
	}
}
