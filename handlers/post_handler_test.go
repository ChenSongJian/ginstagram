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
	"github.com/ChenSongJian/ginstagram/utils"
	"github.com/gin-gonic/gin"
)

func TestListPublicPost_NoPost(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.ListPublicPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":0"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPublicPost_PrivatePost(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	mockPostService.UserService.Users["email"] = models.User{
		Id:        1,
		IsPrivate: true,
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  1,
		Title:   "private post",
		Content: "private post",
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.ListPublicPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":0"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPublicPost_PublicPost(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	mockPostService.UserService.Users["email"] = models.User{
		Id:        1,
		IsPrivate: false,
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  1,
		Title:   "private post",
		Content: "private post",
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.ListPublicPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPublicPost_PostWithMedia(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	mockPostService.UserService.Users["email"] = models.User{
		Id:        1,
		IsPrivate: false,
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  1,
		Title:   "private post",
		Content: "private post",
	}
	mockMediaService.Media[1] = mocks.MediaRecord{
		PostId: 1,
		Url:    "m0",
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.ListPublicPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
	expectedResponseBodyString = "media\":[\"m0"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPublicPost_FilterTitleKeyword(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	mockPostService.UserService.Users["email"] = models.User{
		Id:        1,
		IsPrivate: false,
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  1,
		Title:   "test1",
		Content: "test",
	}
	mockPostService.Posts[2] = mocks.PostRecord{
		UserId:  1,
		Title:   "test",
		Content: "test",
	}

	context.Request, _ = http.NewRequest("GET", "/?keyword=test1", nil)

	handlers.ListPublicPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
	expectedResponseBodyString = "title\":\"test1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPublicPost_FilterContentKeyword(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	mockPostService.UserService.Users["email"] = models.User{
		Id:        1,
		IsPrivate: false,
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  1,
		Title:   "test",
		Content: "test1",
	}
	mockPostService.Posts[2] = mocks.PostRecord{
		UserId:  1,
		Title:   "test",
		Content: "test",
	}

	context.Request, _ = http.NewRequest("GET", "/?keyword=test1", nil)

	handlers.ListPublicPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
	expectedResponseBodyString = "content\":\"test1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPublicPost_Pagination(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	mockPostService.UserService.Users["email"] = models.User{
		Id:        1,
		IsPrivate: false,
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  1,
		Title:   "test",
		Content: "test1",
	}
	mockPostService.Posts[2] = mocks.PostRecord{
		UserId:  1,
		Title:   "test2",
		Content: "test",
	}

	context.Request, _ = http.NewRequest("GET", "/?pageNum=2&pageSize=1", nil)

	handlers.ListPublicPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	var responseBody utils.PageResponse
	json.Unmarshal(response.Body.Bytes(), &responseBody)
	if responseBody.TotalRecords != 2 {
		t.Errorf("Expected total records %d, got %d", 2, responseBody.TotalRecords)
	}
	if len(responseBody.Data.([]interface{})) != 1 {
		t.Errorf("Expected post response length %d", 1)
	}
}

func TestListPost_MissingToken(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.ListPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "user not found in token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}

}

func TestListPost_NoRecord(t *testing.T) {
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

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":0"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}

}

func TestListPost_PublicUserViewOwnPost(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:        1,
		Username:  "test",
		Email:     "test@test.com",
		IsPrivate: false,
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  1,
		Title:   "test",
		Content: "test",
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPost_PrivateUserViewOwnPost(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:        1,
		Username:  "test",
		Email:     "test@test.com",
		IsPrivate: true,
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  1,
		Title:   "test",
		Content: "test",
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPost_UserViewPublicPost(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:        1,
		Username:  "test",
		Email:     "test@test.com",
		IsPrivate: true,
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	mockPostService.UserService.Users["email@email.com"] = models.User{
		Id:        2,
		IsPrivate: false,
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  2,
		Title:   "test",
		Content: "test",
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPost_UserViewPrivatePostWithoutFollowing(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:        1,
		Username:  "test",
		Email:     "test@test.com",
		IsPrivate: true,
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	mockPostService.UserService.Users["email@email.com"] = models.User{
		Id:        2,
		IsPrivate: true,
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  2,
		Title:   "test",
		Content: "test",
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":0"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPost_UserViewPrivatePostWithFollowing(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:        1,
		Username:  "test",
		Email:     "test@test.com",
		IsPrivate: true,
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	mockPostService.UserService.Users["email@email.com"] = models.User{
		Id:        2,
		IsPrivate: true,
	}
	mockPostService.FollowService.Follows[1] = mocks.FollowRecord{
		FolloweeId: 2,
		FollowerId: 1,
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  2,
		Title:   "test",
		Content: "test",
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPost_PostWithMedia(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:        1,
		Username:  "test",
		Email:     "test@test.com",
		IsPrivate: true,
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	mockPostService.UserService.Users["email@email.com"] = models.User{
		Id:        2,
		IsPrivate: true,
	}
	mockPostService.FollowService.Follows[1] = mocks.FollowRecord{
		FolloweeId: 2,
		FollowerId: 1,
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  2,
		Title:   "test",
		Content: "test",
	}
	mockPostService.MediaService.Media[1] = mocks.MediaRecord{
		PostId: 1,
		Url:    "m0",
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
	expectedResponseBodyString = "media\":[\"m0"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPost_SearchTitleKeyword(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:        1,
		Username:  "test",
		Email:     "test@test.com",
		IsPrivate: true,
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  1,
		Title:   "test",
		Content: "test",
	}
	mockPostService.Posts[2] = mocks.PostRecord{
		UserId:  1,
		Title:   "test2",
		Content: "test",
	}

	context.Request, _ = http.NewRequest("GET", "/?keyword=test2", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
	expectedResponseBodyString = "title\":\"test2"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPost_SearchContentKeyword(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:        1,
		Username:  "test",
		Email:     "test@test.com",
		IsPrivate: true,
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  1,
		Title:   "test",
		Content: "test",
	}
	mockPostService.Posts[2] = mocks.PostRecord{
		UserId:  1,
		Title:   "test",
		Content: "test2",
	}

	context.Request, _ = http.NewRequest("GET", "/?keyword=test2", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "total_records\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
	expectedResponseBodyString = "content\":\"test2"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestListPost_Pagination(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	testUser := models.User{
		Id:        1,
		Username:  "test",
		Email:     "test@test.com",
		IsPrivate: true,
	}
	token, err := middlewares.GenerateToken(testUser, true)
	if err != nil {
		t.Errorf("Error generating token: %v", err)
		return
	}
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId:  1,
		Title:   "test",
		Content: "test",
	}
	mockPostService.Posts[2] = mocks.PostRecord{
		UserId:  1,
		Title:   "test",
		Content: "test2",
	}

	context.Request, _ = http.NewRequest("GET", "/?pageNum=2&pageSize=1", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	middlewares.AuthMiddleware()(context)
	handlers.ListPosts(mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	var responseBody utils.PageResponse
	json.Unmarshal(response.Body.Bytes(), &responseBody)
	if responseBody.TotalRecords != 2 {
		t.Errorf("Expected total records %d, got %d", 2, responseBody.TotalRecords)
	}
	if len(responseBody.Data.([]interface{})) != 1 {
		t.Errorf("Expected post response length %d", 1)
	}
}

func TestGetPostById_InvalidPostId(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "invalid_id",
		},
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.GetPostById(mockPostService.UserService, mockPostService.FollowService, mockPostService, mockMediaService)(context)
	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, response.Code)
	}
	expectedResponseBodyString := "invalid post id"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestGetPostById_PostNotFound(t *testing.T) {
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

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.GetPostById(mockPostService.UserService, mockPostService.FollowService, mockPostService, mockMediaService)(context)
	if response.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, response.Code)
	}
	expectedResponseBodyString := "post not found"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}
func TestGetPostById_VisitorViewPrivatePost(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	mockPostService.Posts[1] = mocks.PostRecord{
		UserId: 1,
	}
	mockPostService.UserService.Users["email@email.com"] = models.User{
		Id:        1,
		IsPrivate: true,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.GetPostById(mockPostService.UserService, mockPostService.FollowService, mockPostService, mockMediaService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "post is private, please login and retry again"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestGetPostById_VisitorViewPublicPost(t *testing.T) {
	mockPostService := mocks.NewMockPostService()
	mockMediaService := mockPostService.MediaService

	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	mockPostService.Posts[1] = mocks.PostRecord{
		UserId: 1,
	}
	mockPostService.UserService.Users["email@email.com"] = models.User{
		Id:        1,
		IsPrivate: false,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)

	handlers.GetPostById(mockPostService.UserService, mockPostService.FollowService, mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "\"id\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestGetPostById_UserViewOwnPost(t *testing.T) {
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
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId: 1,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	handlers.GetPostById(mockPostService.UserService, mockPostService.FollowService, mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "\"id\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestGetPostById_UserViewPublicPost(t *testing.T) {
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
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId: 2,
	}
	mockPostService.UserService.Users["email@email.com"] = models.User{
		Id:        2,
		IsPrivate: false,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	handlers.GetPostById(mockPostService.UserService, mockPostService.FollowService, mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "\"id\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestGetPostById_UserViewPrivatePostNotFollowing(t *testing.T) {
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
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId: 2,
	}
	mockPostService.UserService.Users["email@email.com"] = models.User{
		Id:        2,
		IsPrivate: true,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	handlers.GetPostById(mockPostService.UserService, mockPostService.FollowService, mockPostService, mockMediaService)(context)
	if response.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, response.Code)
	}
	expectedResponseBodyString := "post is private and you are not following the author"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestGetPostById_UserViewPrivatePostAndFollowing(t *testing.T) {
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
	mockPostService.Posts[1] = mocks.PostRecord{
		UserId: 2,
	}
	mockPostService.UserService.Users["email@email.com"] = models.User{
		Id:        2,
		IsPrivate: false,
	}
	mockPostService.FollowService.Follows[1] = mocks.FollowRecord{
		FollowerId: 1,
		FolloweeId: 2,
	}

	context.Params = []gin.Param{
		{
			Key:   "postId",
			Value: "1",
		},
	}

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer "+token)

	handlers.GetPostById(mockPostService.UserService, mockPostService.FollowService, mockPostService, mockMediaService)(context)
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}
	expectedResponseBodyString := "\"id\":1"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

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
