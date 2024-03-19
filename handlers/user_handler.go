package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ChenSongJian/ginstagram/models"
	"github.com/ChenSongJian/ginstagram/services"
	"github.com/ChenSongJian/ginstagram/utils"
	"github.com/gin-gonic/gin"
)

type UserRegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Id              int    `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Bio             string `json:"bio"`
	ProfileImageUrl string `json:"profile_image_url"`
}

func RegisterUser(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UserRegisterReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if !utils.IsComplex(req.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be complex"})
			return
		}

		passwordHash := utils.GenerateHash(req.Password)
		user := models.User{
			Username:     req.Username,
			PasswordHash: passwordHash,
			Email:        req.Email,
		}

		if err := userService.Create(user); err != nil {
			duplicateErrorMsg := "ERROR: duplicate key value violates unique constraint"
			if strings.Contains(err.Error(), duplicateErrorMsg) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}
}

func GetUserById(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.Param("userId")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		user, err := userService.GetById(userId)
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userResponse := UserResponse{
			Id:              user.Id,
			Username:        user.Username,
			Email:           user.Email,
			Bio:             user.Bio,
			ProfileImageUrl: user.ProfileImageUrl,
		}
		c.JSON(http.StatusOK, userResponse)
	}
}
