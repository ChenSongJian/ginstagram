package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ChenSongJian/ginstagram/middlewares"
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

type UserLoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	Id              int    `json:"id"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Bio             string `json:"bio"`
	ProfileImageUrl string `json:"profile_image_url"`
	IsPrivate       bool   `json:"is_private"`
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

func ListUsers(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageNum := c.Query("pageNum")
		pageSize := c.Query("pageSize")
		keyword := c.Query("keyword")
		users, pageInfo, err := userService.List(pageNum, pageSize, keyword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userResponses := make([]UserResponse, len(users))
		for i, user := range users {
			userResponses[i] = UserResponse{
				Id:              user.Id,
				Username:        user.Username,
				Email:           user.Email,
				Bio:             user.Bio,
				ProfileImageUrl: user.ProfileImageUrl,
				IsPrivate:       user.IsPrivate,
			}
		}
		pageInfo.Data = userResponses
		c.JSON(http.StatusOK, pageInfo)
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
			IsPrivate:       user.IsPrivate,
		}
		c.JSON(http.StatusOK, userResponse)
	}
}

func LoginUser(userService services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UserLoginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := userService.GetByEmail(req.Email)
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !utils.CompareHash(user.PasswordHash, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
			return
		}
		token, err := middlewares.GenerateToken(user, true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
