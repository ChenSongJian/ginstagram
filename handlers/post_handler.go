package handlers

import (
	"net/http"
	"strconv"

	"github.com/ChenSongJian/ginstagram/models"
	"github.com/ChenSongJian/ginstagram/services"
	"github.com/gin-gonic/gin"
)

type PostReq struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Media   []string `json:"media"`
}

type PostResponse struct {
	Id        int      `json:"id"`
	CreatedAt string   `json:"created_at"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	UserId    int      `json:"user_id"`
	Media     []string `json:"media"`
}

func CreatePost(postService services.PostService, mediaService services.MediaService) gin.HandlerFunc {
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
		var req PostReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		post := models.Post{
			Title:   req.Title,
			Content: req.Content,
			UserId:  modelTokenUser.Id,
		}
		if len(req.Media) < 1 || len(req.Media) > 9 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "please upload at least one and no more than 9 media"})
			return
		}
		postId, err := postService.Create(post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var media []models.Media
		for _, m := range req.Media {
			media = append(media, models.Media{
				Url:    m,
				PostId: postId,
			})
		}
		if err := mediaService.Create(media); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			_ = postService.DeleteById(postId)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "post created successfully"})
	}
}

func DeletePost(postService services.PostService, mediaService services.MediaService) gin.HandlerFunc {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if post.UserId != modelTokenUser.Id {
			c.JSON(http.StatusForbidden, gin.H{"error": "no permission to delete this post"})
			return
		}
		err = postService.DeleteById(modelTokenUser.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
	}
}
