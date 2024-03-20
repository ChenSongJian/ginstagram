package web

import (
	"github.com/ChenSongJian/ginstagram/handlers"
	"github.com/ChenSongJian/ginstagram/middlewares"
	"github.com/ChenSongJian/ginstagram/services"
	"github.com/gin-gonic/gin"
)

var userService services.UserService
var followService services.FollowService
var postService services.PostService
var mediaService services.MediaService
var commentService services.CommentService

func initServices() {
	userService = services.NewDBUserService()
	followService = services.NewDBFollowService()
	postService = services.NewDBPostService()
	mediaService = services.NewDBMediaService()
	commentService = services.NewDBCommentService()
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	initServices()

	apiV1Group := r.Group("/api/v1")
	userV1Group := apiV1Group.Group("/user")
	userV1Group.POST("/", handlers.RegisterUser(userService))
	userV1Group.GET("/", handlers.ListUsers(userService))
	userV1Group.GET("/:userId", handlers.GetUserById(userService))
	userV1Group.PUT("/:userId", middlewares.AuthMiddleware(), handlers.UpdateUser(userService))
	userV1Group.DELETE("/:userId", middlewares.AuthMiddleware(), handlers.DeleteUser(userService))
	userV1Group.POST("/login", handlers.LoginUser(userService))
	userV1Group.POST("/logout", middlewares.AuthMiddleware(), handlers.LogoutUser(userService))
	userV1Group.GET("/refresh", middlewares.AuthMiddleware(), handlers.RefreshToken(userService))

	followV1Group := apiV1Group.Group("/follow")
	followV1Group.GET("/", handlers.ListFollows(followService))
	followV1Group.POST("/", middlewares.AuthMiddleware(), handlers.FollowUser(followService))
	followV1Group.DELETE("/:followId", middlewares.AuthMiddleware(), handlers.UnfollowUser(followService))

	postV1Group := apiV1Group.Group("/post")
	postV1Group.GET("/public", handlers.ListPublicPosts(postService, mediaService))
	postV1Group.GET("/", middlewares.AuthMiddleware(), handlers.ListPosts(postService, mediaService))
	postV1Group.POST("/", middlewares.AuthMiddleware(), handlers.CreatePost(postService, mediaService))
	postV1Group.DELETE("/:postId", middlewares.AuthMiddleware(), handlers.DeletePost(postService, mediaService))
	postV1Group.GET("/:postId/comment", middlewares.AuthMiddleware(), handlers.ListCommentsByPostId(userService, followService, postService, commentService))
	postV1Group.POST("/:postId/comment", middlewares.AuthMiddleware(), handlers.CreateComment(userService, followService, postService, commentService))
	postV1Group.DELETE("/:postId/comment/:commentId", middlewares.AuthMiddleware(), handlers.DeleteComment(userService, followService, postService, commentService))
	return r
}
