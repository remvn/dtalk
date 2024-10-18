package middleware

import (
	"dtalk/internal/app/dtalk"
	"dtalk/internal/config"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Auth struct {
	tokenConfig config.JwtTokenConfig
}

func NewAuth(tokenConfig config.JwtTokenConfig) *Auth {
	return &Auth{
		tokenConfig: tokenConfig,
	}
}

func (m *Auth) Apply(consumer consumer) {
	consumer.Use(m.validate(), m.extract())
}

func (m *Auth) validate() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContextKey:  UserTokenContextKey,
		SigningKey:  m.tokenConfig.Secret,
		TokenLookup: fmt.Sprintf("header:Authorization:Bearer ,cookie:%s", m.tokenConfig.Name),
	})
}

var ErrUnableToExtract = errors.New("unable to extract data from context")

func (m *Auth) extract() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get(UserTokenContextKey).(*jwt.Token)
			if !ok {
				log.Println(fmt.Errorf("unable to assert *jwt.Token: %w", ErrUnableToExtract))
				return c.NoContent(http.StatusInternalServerError)
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Println(fmt.Errorf("unable to assert jwt.MapClaims: %w", ErrUnableToExtract))
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

func (m *Auth) Name() string {
	return "Auth"
}
