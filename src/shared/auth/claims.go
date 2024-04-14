package auth

import (
	"errors"
	"fmt"
	"github.com/JECSand/eventit-server/src/shared/enums"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type AppClaims struct {
	ProfileId string     `json:"profileId,omitempty"`
	Role      enums.Role `json:"role,omitempty"`
	jwt.RegisteredClaims
}

func DecodeJWT(tokenString string) (*AppClaims, error) {
	if tokenString == "" {
		return &AppClaims{}, errors.New("unauthorized")
	}
	secret := viper.GetString("auth_jwt_secret")
	token, err := jwt.ParseWithClaims(tokenString, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		//log.Fatal(err)
		return &AppClaims{}, err
	} else if claims, ok := token.Claims.(*AppClaims); ok && token.Valid {
		return claims, nil
	}
	return &AppClaims{}, errors.New("unknown claims type, cannot proceed")
}
