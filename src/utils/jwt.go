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

func ParseToken(tokenStr string) (*JwtCustomClaims, error) {
	claims := &JwtCustomClaims{}

	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return signingKey, nil
	})

	// 如果解析过程中出现错误，处理不同错误类型
	if err != nil {
		// 签名无效
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid signature")
		}

		// Token 已过期
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}

		// 其他解析错误
		return nil, err
	}

	// 检查 Token 是否有效
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func RefreshToken(tokenStr string) (string, error) {
	// 解析旧的 Token
	claims, err := ParseToken(tokenStr)
	if err != nil {
		return "", err
	}

	// 更新过期时间
	// todo 续签可以适当延长
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.tokenExpire") * time.Minute))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	// 创建新的 Token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 返回新的 Token 字符串
	return newToken.SignedString(signingKey)
}