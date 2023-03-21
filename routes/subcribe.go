package routes

import (
	"wayshub/handlers"
	"wayshub/pkg/middleware"
	"wayshub/pkg/mysql"
	"wayshub/repositories"

	"github.com/labstack/echo/v4"
)

func SubscribeRoutes(e *echo.Echo) {
	subscribeRepository := repositories.RepositorySubscribe(mysql.DB)
	h := handlers.HandlerSubscribe(subscribeRepository)

	e.GET("/subscribes", h.FindSubscribes, middleware.Auth)
	e.GET("/subscribe/:id", h.GetSubscribe, middleware.Auth)

	e.GET("/subscribeByOther/:id", h.GetSubscribeByOther, middleware.Auth)

	e.POST("/subscribe/:id", h.CreateSubscribe, middleware.Auth)
	e.DELETE("/subscribe", h.DeleteSubscribe, middleware.Auth)

	e.GET("/subscription", h.GetSubscription, middleware.Auth)
}
