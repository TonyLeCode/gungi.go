package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func VerifySupabaseTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		splitToken := strings.Split(c.Request().Header.Get("Authorization"), "Bearer")
		if len(splitToken) != 2 {
			return c.JSON(http.StatusUnauthorized, "Invalid Token")
		}

		tokenString := strings.TrimSpace(splitToken[1])

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			viper.SetConfigFile("../.env")
			viper.ReadInConfig()

			supabaseSecret := viper.GetString("SUPABASE_JWT_SECRET")

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(supabaseSecret), nil
		})
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusUnauthorized, "Invalid Token")
		}

		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Invalid Token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, "Invalid Token")
		}

		c.Set("sub", claims["sub"])

		return next(c)
	}
}
