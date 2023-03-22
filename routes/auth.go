package routes

import (
	"wayshub/handlers"
	"wayshub/pkg/middleware"
	"wayshub/pkg/mysql"
	"wayshub/repositories"

	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Group) {
	channelRepository := repositories.RepositoryChannel(mysql.DB)
	h := handlers.HandlerAuth(channelRepository)

	e.POST("/register", h.Register)
	e.POST("/login", h.Login)
	e.GET("/check-auth", h.CheckAuth, middleware.Auth)
}
