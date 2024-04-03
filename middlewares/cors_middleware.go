package middlewares

import (
	"github.com/gin-contrib/cors"
)

func NewCorsConfig() cors.Config {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	return config
}
