package middleware

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Cors 跨域配置
func Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Cookie"}
	config.AllowCredentials = true

	origins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	config.AllowOrigins = make([]string, 0, len(origins))
	for _, origin := range origins {
		if origin = strings.TrimSpace(origin); origin != "" {
			config.AllowOrigins = append(config.AllowOrigins, origin)
		}
	}
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = []string{"http://localhost:8080"}
	}
	return cors.New(config)
}
