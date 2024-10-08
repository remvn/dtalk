package middleware

import (
	"dtalk/internal/config"
	"dtalk/internal/dtalk"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	tokenConfig config.JwtTokenConfig
}

func NewAuthMiddleware(tokenConfig config.JwtTokenConfig) *AuthMiddleware {
	return &AuthMiddleware{
		tokenConfig: tokenConfig,
	}
}

func (m *AuthMiddleware) Apply(consumer consumer) {
	consumer.Use(m.validate(), m.extract())
}

func (m *AuthMiddleware) validate() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey:  UserClaimsContextKey,
		SigningKey:  m.tokenConfig.Secret,
		TokenLookup: fmt.Sprintf("header:Authorization:Bearer ,cookie:%s", m.tokenConfig.Name),
	})
}

var ErrUnableToExtract = errors.New("unable to extract user claims from context")

func (m *AuthMiddleware) extract() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, ok := c.Get(UserClaimsContextKey).(jwt.MapClaims)
			if !ok {
				log.Println(ErrUnableToExtract)
				return c.NoContent(http.StatusInternalServerError)
			}

			// TODO: consider parsing with https://github.com/go-viper/mapstructure ?
			name, ok := claims["name"].(string)
			if !ok {
				log.Println(fmt.Errorf("%s field is undefined: %w", "name", ErrUnableToExtract))
				return c.NoContent(http.StatusInternalServerError)
			}
			id, ok := claims["id"].(string)
			if !ok {
				log.Println(fmt.Errorf("%s field is undefined: %w", "id", ErrUnableToExtract))
				return c.NoContent(http.StatusInternalServerError)
			}

			userInfo := &dtalk.UserTokenInfo{
				Name: name,
				ID:   id,
			}
			c.Set(UserInfoContextKey, userInfo)

			return next(c)
		}
	}
}
