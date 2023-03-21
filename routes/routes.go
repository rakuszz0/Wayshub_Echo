package routes

import "github.com/labstack/echo/v4"

func RouteInit(e *echo.Group) {
	AuthRoutes(e)
	ChannelRoutes(e)
	VideoRoutes(e)
	CommentRoutes(e)
	SubscribeRoutes(e)
}
