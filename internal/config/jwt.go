package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenConfig struct {
	Name     string
	Secret   []byte
	Duration time.Duration
}

func (config JwtTokenConfig) Sign(claims jwt.MapClaims) (string, error) {
	claims["iat"] = time.Now().Unix()
	claims["exp"] = config.GetExpire().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(config.Secret)
}

func (config JwtTokenConfig) Parse(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// by default this library will prevent token with alg=none
		// https://pkg.go.dev/github.com/golang-jwt/jwt/v5#readme-compliance
		// this is a double check just to make sure
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return config.Secret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("invalid jwt claims type")
	}
}

func (config *JwtTokenConfig) GetExpire() time.Time {
	return time.Now().Add(config.Duration)
}
