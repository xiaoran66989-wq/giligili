package server

import (
	"giligili/api"
	"giligili/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())

	v1 := r.Group("/api/v1")
	{
		v1.POST("ping", api.Ping)
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)

		authed := v1.Group("/")
		authed.Use(middleware.AuthRequired())
		{
			authed.GET("user/me", api.UserMe)
			authed.DELETE("user/logout", api.UserLogout)
			authed.POST("videos", api.CreateVideo)
			authed.PUT("video/:id", api.UpdateVideo)
			authed.DELETE("video/:id", api.DeleteVideo)
			authed.POST("upload/token", api.UploadToken)
		}

		v1.GET("video/:id", api.ShowVideo)
		v1.GET("videos", api.ListVideo)
	}

	r.StaticFile("/swagger.json", "./swagger/swagger.json")
	r.Static("/swagger", "./swagger/dist")
	return r
}
