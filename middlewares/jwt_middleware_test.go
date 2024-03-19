package middlewares_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ChenSongJian/ginstagram/middlewares"
	"github.com/ChenSongJian/ginstagram/models"
	"github.com/gin-gonic/gin"
)

func TestGenerateToken_MissingUser(t *testing.T) {
	token, err := middlewares.GenerateToken(models.User{}, true)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if token != "" {
		t.Error("Expected empty token, got", token)
	}
}

func TestGenerateToken_ValidUser(t *testing.T) {
	user := models.User{
		Id:           1,
		Username:     "username",
		PasswordHash: "PasswordHash",
		Email:        "email",
	}
	token, err := middlewares.GenerateToken(user, true)
	if err != nil {
		t.Error("Expected no error, got", err)
	}
	if token == "" {
		t.Error("Expected non-empty token, got empty")
	}
}

func TestAuthMiddleware_NoAuthorizationHeader(t *testing.T) {
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request, _ = http.NewRequest("GET", "/", nil)
	middlewares.AuthMiddleware()(context)
	if response.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, response.Code)
	}
	expectedResponseBodyString := "Authorization header is required"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body to contain %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestAuthMiddleware_InvalidAuthorizationFormat(t *testing.T) {
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "token")
	middlewares.AuthMiddleware()(context)
	if response.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, response.Code)
	}
	expectedResponseBodyString := "Invalid Authorization header"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body to contain %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", "Bearer invalid_token")
	middlewares.AuthMiddleware()(context)
	if response.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, response.Code)
	}
	expectedResponseBodyString := "Failed to parse token"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body to contain %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestAuthMiddleware_InactiveValidToken(t *testing.T) {
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	user := models.User{
		Id:           1,
		Username:     "username",
		PasswordHash: "PasswordHash",
		Email:        "email",
	}
	token, _ := middlewares.GenerateToken(user, false)
	formatted_token := fmt.Sprintf("Bearer %s", token)

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", formatted_token)
	middlewares.AuthMiddleware()(context)

	if response.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, response.Code)
	}
	expectedResponseBodyString := "Token is not active"
	if !strings.Contains(response.Body.String(), expectedResponseBodyString) {
		t.Errorf("Expected response body to contain %s, got %s", expectedResponseBodyString, response.Body.String())
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)

	user := models.User{
		Id:           1,
		Username:     "username",
		PasswordHash: "PasswordHash",
		Email:        "email",
	}
	token, _ := middlewares.GenerateToken(user, true)
	formatted_token := fmt.Sprintf("Bearer %s", token)

	context.Request, _ = http.NewRequest("GET", "/", nil)
	context.Request.Header.Set("Authorization", formatted_token)
	middlewares.AuthMiddleware()(context)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, response.Code)
	}

	tokenUser := context.MustGet("user").(models.User)
	if tokenUser.Id != user.Id {
		t.Errorf("Expected user ID %d, got %d", user.Id, tokenUser.Id)
	}
	if tokenUser.Username != user.Username {
		t.Errorf("Expected user username %s, got %s", user.Username, tokenUser.Username)
	}
	if tokenUser.PasswordHash != user.PasswordHash {
		t.Errorf("Expected user password hash %s, got %s", user.PasswordHash, tokenUser.PasswordHash)
	}
	if tokenUser.Email != user.Email {
		t.Errorf("Expected user email %s, got %s", user.Email, tokenUser.Email)
	}
	if tokenUser.Bio != user.Bio {
		t.Errorf("Expected user bio %s, got %s", user.Bio, tokenUser.Bio)
	}
	if tokenUser.ProfileImageUrl != user.ProfileImageUrl {
		t.Errorf("Expected user profile image URL %s, got %s", user.ProfileImageUrl, tokenUser.ProfileImageUrl)
	}
}
