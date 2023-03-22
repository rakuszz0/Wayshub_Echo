package routes

import (
	"wayshub/handlers"
	"wayshub/pkg/middleware"
	"wayshub/pkg/mysql"
	"wayshub/repositories"

	"github.com/labstack/echo/v4"
)

func ChannelRoutes(e *echo.Group) {
	channelRepository := repositories.RepositoryChannel(mysql.DB)
	h := handlers.HandlerChannel(channelRepository)

	e.GET("/channels", h.FindChannels)
	e.GET("/channel/:id", h.GetChannel)
	e.PATCH("/channel/:id", h.UpdateChannel, middleware.Auth, middleware.UploadCover, middleware.UploadPhoto)
	e.DELETE("/channel/:id", h.DeleteChannel, middleware.Auth)

	e.PATCH("/plusSubs/:id", h.PlusSubscriber, middleware.Auth)
	e.PATCH("/minusSubs/:id", h.MinusSubscriber, middleware.Auth)
}
