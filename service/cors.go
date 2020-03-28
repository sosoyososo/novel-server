package service

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func corsMiddleWare() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	config.AllowMethods = []string{
		"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
	}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Token",
		"x-requested-with",
	}
	config.AllowOrigins = []string{
		"https://inhooer.com",
		"https://h5.inhooer.com",
		"https://dev.inhooer.com",

		"http://inhooer.com",
		"http://h5.inhooer.com",
		"http://dev.inhooer.com",

		"http://localhost",
		"http://127.0.0.1",
		"http://0.0.0.0",
	}
	config.AllowOriginFunc = func(origin string) bool {
		return strings.HasPrefix(origin, "http://localhost") ||
			strings.HasPrefix(origin, "http://127.0.0.1")
	}
	return cors.New(config)
}
