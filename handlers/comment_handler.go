package handlers

import (
	"net/http"
	"strconv"

	"github.com/ChenSongJian/ginstagram/models"
	"github.com/ChenSongJian/ginstagram/services"
	"github.com/gin-gonic/gin"
)

type CommentReq struct {
	Content string `json:"content" binding:"required"`
}

func CreateComment(userService services.UserService, followService services.FollowService,
	postService services.PostService, commentService services.CommentService) gin.HandlerFunc {
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
				c.JSON(http.StatusBadRequest, gin.H{"error": "post not found"})
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
		var req CommentReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := commentService.Create(postId, modelTokenUser.Id, req.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "comment created successfully"})
	}
}
