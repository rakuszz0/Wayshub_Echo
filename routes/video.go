package routes

import (
	"wayshub/handlers"
	"wayshub/pkg/middleware"
	"wayshub/pkg/mysql"
	"wayshub/repositories"

	"github.com/labstack/echo/v4"
)

func VideoRoutes(e *echo.Echo) {
	videoRepository := repositories.RepositoryVideo(mysql.DB)
	h := handlers.HandlerVideo(videoRepository)

	e.GET("/videos", h.FindVideos)
	e.GET("/video/:id", h.GetVideo)

	e.POST("/video", h.CreateVideo, middleware.Auth, middleware.UploadVideo, middleware.UploadThumbnail)
	e.PATCH("/video/:id", h.UpdateVideo, middleware.Auth, middleware.UploadVideo, middleware.UploadThumbnail)

	e.DELETE("/video/:id", h.DeleteVideo, middleware.Auth)

	e.GET("/myvideo", h.FindVideosByChannelId, middleware.Auth)
	e.GET("/FindMyVideos", h.FindMyVideos, middleware.Auth)

	e.PATCH("/UpdateViews/:id", h.UpdateViews, middleware.Auth)
}
