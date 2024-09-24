package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"time"
)

var signingKey = []byte(viper.GetString("jwt.signingKey"))

type JwtCustomClaims struct {
	ID   uint
	Name string
	jwt.RegisteredClaims
}

func GenerateToken(id uint, name string) (string, error) {
	jwtCustomClaims := JwtCustomClaims{
		ID:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.tokenExpire") * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "Token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtCustomClaims)
	return token.SignedString(signingKey)
}

func ParseToken(tokenStr string) (JwtCustomClaims, error) {
	jwtCustomClaims := JwtCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &jwtCustomClaims, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err == nil && !token.Valid {
		err = errors.New("token is invalid")
	}
	return jwtCustomClaims, err
}

func IsTokenValid(tokenStr string) bool {
	_, err := ParseToken(tokenStr)
	if err != nil {
		return false
	}
	return true
}
