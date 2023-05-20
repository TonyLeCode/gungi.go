package middleware

import (
	"net/http"

	"github.com/TonyLeCode/gungi.go/server/auth"
	"github.com/labstack/echo/v4"
)

func VerifySupabaseTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := auth.AuthenticateSupabaseToken(c.Request().Header.Get("Authorization"))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err)
		}

		c.Set("sub", claims["sub"])

		return next(c)
	}
}
