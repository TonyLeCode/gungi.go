package auth

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func AuthenticateSupabaseToken(tokenBearer string) (jwt.MapClaims, error) {
	splitToken := strings.Split(tokenBearer, "Bearer")
	if len(splitToken) != 2 {
		return nil, errors.New("invalid token")
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
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
