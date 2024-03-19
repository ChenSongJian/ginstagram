package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ChenSongJian/ginstagram/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = "my-jwt-key"

func GenerateToken(user models.User, active bool) (string, error) {
	if user == (models.User{}) {
		return "", errors.New("invalid user")
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["active"] = active
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			return
		}

		token, err := jwt.Parse(splitToken[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Failed to extract claims"})
			return
		}

		expiry := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expiry) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			return
		}

		active, ok := claims["active"].(bool)
		if !ok || !active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is not active"})
			return
		}

		user, err := extractUserFromClaims(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Failed to extract user information: %s", err)})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func extractUserFromClaims(claims jwt.MapClaims) (models.User, error) {
	user := models.User{}
	claimsUser, ok := claims["user"].(map[string]interface{})
	if !ok {
		return user, errors.New("failed to extract user information")
	}
	fmt.Println(claimsUser)

	userId, ok := claimsUser["Id"].(float64)
	if !ok {
		return user, errors.New("failed to extract user Id")
	}
	user.Id = int(userId)
	user.Username, ok = claimsUser["Username"].(string)
	if !ok {
		return user, errors.New("failed to extract user username")
	}
	user.PasswordHash, ok = claimsUser["PasswordHash"].(string)
	if !ok {
		return user, errors.New("failed to extract user password hash")
	}
	user.Email, ok = claimsUser["Email"].(string)
	if !ok {
		return user, errors.New("failed to extract user email")
	}
	user.Bio, ok = claimsUser["Bio"].(string)
	if !ok {
		return user, errors.New("failed to extract user bio")
	}
	user.ProfileImageUrl, ok = claimsUser["ProfileImageUrl"].(string)
	if !ok {
		return user, errors.New("failed to extract user profile image url")
	}
	user.IsPrivate, ok = claimsUser["IsPrivate"].(bool)
	if !ok {
		return user, errors.New("failed to extract user is private")
	}
	return user, nil
}
