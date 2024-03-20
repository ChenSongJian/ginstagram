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

type FollowResponse struct {
	Id         int `json:"id"`
	FollowerId int `json:"follower_id"`
	FolloweeId int `json:"followee_id"`
}

func ListFollows(followService services.FollowService) gin.HandlerFunc {
	return func(c *gin.Context) {
		followerIdStr := c.Query("follower_id")
		followeeIdStr := c.Query("followee_id")
		if (followerIdStr == "" && followeeIdStr == "") || (followerIdStr != "" && followeeIdStr != "") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Either provide follower_id or followee_id, but not both and not neither."})
			return
		}
		var follows []models.Follow
		if followerIdStr != "" {
			var followerId int
			followerId, err := strconv.Atoi(followerIdStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid follower_id"})
				return
			}
			follows, err = followService.GetByFollowerId(followerId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		if followeeIdStr != "" {
			var followeeId int
			followeeId, err := strconv.Atoi(followeeIdStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid followee_id"})
				return
			}
			follows, err = followService.GetByFolloweeId(followeeId)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"follows": follows})
	}
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
