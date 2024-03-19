package web

import (
	"github.com/ChenSongJian/ginstagram/handlers"
	"github.com/ChenSongJian/ginstagram/services"
	"github.com/gin-gonic/gin"
)

var userService services.UserService

func initServices() {
	userService = services.NewDBUserService()
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
	userV1Group.GET("/:userId", handlers.GetUserById(userService))
	userV1Group.POST("/", handlers.RegisterUser(userService))

	return r
}
