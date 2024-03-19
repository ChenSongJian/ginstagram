package handlers

import (
	"net/http"

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
			if err.Error() == "ERROR: insert or update on table \"follows\" violates foreign key constraint \"fk_user\"" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
				return
			}
			if err.Error() == "ERROR: new row for relation \"follows\" violates check constraint \"different_user_and_follower\"" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "can not follow yourself"})
				return
			}
			if err.Error() == "ERROR: duplicate key value violates unique constraint \"unique_user_follower_pair\"" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "already following"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "follow user success"})
	}
}
