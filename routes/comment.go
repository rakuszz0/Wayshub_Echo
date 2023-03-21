package routes

import (
	"wayshub/handlers"
	"wayshub/pkg/middleware"
	"wayshub/pkg/mysql"
	"wayshub/repositories"

	"github.com/labstack/echo/v4"
)

func CommentRoutes(e *echo.Echo) {
	commentRepository := repositories.RepositoryComment(mysql.DB)
	h := handlers.HandlerComment(commentRepository)

	e.GET("/comments", h.FindComments)
	e.GET("/comment/:id", h.GetComment)
	e.POST("/comment/:id", h.CreateComment, middleware.Auth)
	e.PATCH("/comment/:id", h.UpdateComment, middleware.Auth)
	e.DELETE("/comment/:id", h.DeleteComment, middleware.Auth)
}
