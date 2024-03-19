package web

import (
	"github.com/ChenSongJian/ginstagram/handlers"
	"github.com/ChenSongJian/ginstagram/middlewares"
	"github.com/ChenSongJian/ginstagram/services"
	"github.com/gin-gonic/gin"
)

var userService services.UserService
var followService services.FollowService

func initServices() {
	userService = services.NewDBUserService()
	followService = services.NewDBFollowService()
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
	followV1Group.POST("/", middlewares.AuthMiddleware(), handlers.FollowUser(followService))

	return r
}
