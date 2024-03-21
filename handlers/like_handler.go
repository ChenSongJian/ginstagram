package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ChenSongJian/ginstagram/models"
	"github.com/ChenSongJian/ginstagram/services"
	"github.com/gin-gonic/gin"
)

type LikeResponse struct {
	Id     int `json:"id"`
	UserId int `json:"user_id"`
}

func ListLikesByPostId(userService services.UserService, followService services.FollowService,
	postService services.PostService, likeService services.LikeService) gin.HandlerFunc {
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
		postIdStr := c.Param("postId")
		postId, err := strconv.Atoi(postIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
			return
		}
		var post models.Post
		post, err = postService.GetById(postId)
		if err != nil {
			if err.Error() == "record not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if post.UserId != modelTokenUser.Id {
			var author models.User
			author, _ = userService.GetById(post.UserId)
			if author.IsPrivate {
				if !followService.IsFollowing(modelTokenUser.Id, author.Id) {
					c.JSON(http.StatusForbidden, gin.H{"error": "post is private and you are not following the author"})
					return
				}
			}
		}
		var likes []models.PostLike
		likes, err = likeService.ListPostLikesByPostId(postId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		likeResponses := make([]LikeResponse, len(likes))
		for i, like := range likes {
			likeResponses[i] = LikeResponse{
				Id:     like.Id,
				UserId: like.UserId,
			}
		}
		c.JSON(http.StatusOK, gin.H{"likes": likeResponses})
	}
}

func LikePost(userService services.UserService, followService services.FollowService,
	postService services.PostService, likeService services.LikeService) gin.HandlerFunc {
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
		postIdStr := c.Param("postId")
		postId, err := strconv.Atoi(postIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
			return
		}
		var post models.Post
		post, err = postService.GetById(postId)
		if err != nil {
			if err.Error() == "record not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if post.UserId != modelTokenUser.Id {
			var author models.User
			author, _ = userService.GetById(post.UserId)
			if author.IsPrivate {
				if !followService.IsFollowing(modelTokenUser.Id, author.Id) {
					c.JSON(http.StatusForbidden, gin.H{"error": "post is private and you are not following the author"})
					return
				}
			}
		}
		if err := likeService.CreatePostLike(postId, modelTokenUser.Id); err != nil {
			if strings.Contains(err.Error(), "violates unique constraint \"unique_post_user_pair\"") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "already liked"})
				return
			}
			if strings.Contains(err.Error(), "violates foreign key constraint \"post_likes_user_id_fkey\"") {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "like created successfully"})

	}
}

func UnlikePost(userService services.UserService, followService services.FollowService,
	postService services.PostService, likeService services.LikeService) gin.HandlerFunc {
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
		postLikeIdStr := c.Param("postLikeId")
		postLikeId, err := strconv.Atoi(postLikeIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post like id"})
			return
		}
		var postLike models.PostLike
		postLike, err = likeService.GetByPostLikeId(postLikeId)
		if err != nil {
			if err.Error() == "record not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "post like not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if postLike.UserId != modelTokenUser.Id {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission to unlike"})
			return
		}
		if err := likeService.DeletePostLikeById(postLikeId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "like deleted successfully"})
	}
}
