package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ChenSongJian/ginstagram/models"
	"github.com/ChenSongJian/ginstagram/services"
	"github.com/gin-gonic/gin"
)

type FollowUserReq struct {
	UserId int `json:"user_id" binding:"required"`
}

func FollowUser(followService services.FollowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenUser, exists := c.Get("tokenUser")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found in token"})
			return
		}
		modelTokenUser, ok := tokenUser.(models.User)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token user type"})
			return
		}
		var req FollowUserReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.UserId == modelTokenUser.Id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "can not follow yourself"})
			return
		}
		if err := followService.Create(modelTokenUser.Id, req.UserId); err != nil {
			if strings.Contains(err.Error(), "violates foreign key constraint \"fk_user\"") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
				return
			}
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "already following"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "follow user success"})
	}
}

func UnfollowUser(followService services.FollowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenUser, exists := c.Get("tokenUser")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found in token"})
			return
		}
		modelTokenUser, ok := tokenUser.(models.User)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token user type"})
			return
		}
		followIdStr := c.Param("followId")
		followId, err := strconv.Atoi(followIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid follow id"})
			return
		}
		follow, err := followService.GetById(followId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "follow not found"})
			return
		}
		if follow.FollowerId != modelTokenUser.Id {
			c.JSON(http.StatusBadRequest, gin.H{"error": "not follower"})
			return
		}

		if err := followService.Delete(followId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "unfollow user success"})
	}
}
