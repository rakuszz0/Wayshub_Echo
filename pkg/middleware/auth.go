package middleware

import (
	// "context"
	// "encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	dto "wayshub/dto/result"
	jwtToken "wayshub/pkg/jwt"
)

type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		if token == "" {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Code: http.StatusBadRequest, Message: "unauthorized"})
		}

		token = strings.Split(token, " ")[1]
		claims, err := jwtToken.DecodeToken(token)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, Result{Code: http.StatusBadRequest, Message: "unathorized"})
		}

		c.Set("channelInfo", claims)
		return next(c)
	}
}
