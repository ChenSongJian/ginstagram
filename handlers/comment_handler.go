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

type CommentResponse struct {
	Id        int    `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UserId    int    `json:"userId"`
}

func ListCommentsByPostId(userService services.UserService, followService services.FollowService,
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
		pageNum := c.Query("pageNum")
		pageSize := c.Query("pageSize")
		comment, pageInfo, err := commentService.ListByPostId(postId, pageNum, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		comments := make([]CommentResponse, 0)
		for _, comment := range comment {
			var user models.User
			user, _ = userService.GetById(comment.UserId)
			comments = append(comments, CommentResponse{
				Id:        comment.Id,
				Content:   comment.Content,
				CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
				UserId:    user.Id,
			})
		}
		pageInfo.Data = comments
		c.JSON(http.StatusOK, pageInfo)
	}
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

func DeleteComment(userService services.UserService, followService services.FollowService,
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
		_, err := strconv.Atoi(postIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
			return
		}
		commentIdStr := c.Param("commentId")
		commentId, err := strconv.Atoi(commentIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
			return
		}
		comment, err := commentService.GetById(commentId)
		if err != nil {
			if err.Error() == "record not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if comment.UserId != modelTokenUser.Id {
			c.JSON(http.StatusForbidden, gin.H{"error": "you are not the author of the comment"})
			return
		}
		if err := commentService.DeleteById(commentId); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
	}
}
