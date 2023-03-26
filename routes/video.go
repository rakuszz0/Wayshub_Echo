package routes

import (
	"wayshub/handlers"
	"wayshub/pkg/middleware"
	"wayshub/pkg/mysql"
	"wayshub/repositories"

	"github.com/labstack/echo/v4"
)

func VideoRoutes(e *echo.Group) {
	videoRepository := repositories.RepositoryVideo(mysql.DB)
	h := handlers.HandlerVideo(videoRepository)

	e.GET("/videos", h.FindVideos)
	e.GET("/video/:id", h.GetVideo)

	e.POST("/video", middleware.Auth(middleware.UploadVideo(middleware.UploadThumbnail(h.CreateVideo))))
	e.PATCH("/video/:id", middleware.Auth(middleware.UploadVideo(middleware.UploadThumbnail(h.UpdateVideo))))

	e.DELETE("/video/:id", middleware.Auth(h.DeleteVideo))

	e.GET("/myvideo", middleware.Auth(h.FindVideosByChannelId))
	e.PATCH("/UpdateViews/:id", middleware.Auth(h.UpdateViews))
}
