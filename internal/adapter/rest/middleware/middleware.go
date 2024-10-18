package middleware

import (
	"dtalk/internal/app/dtalk"
	"errors"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

const (
	UserTokenContextKey = "user_claims"
	UserInfoContextKey  = "user_info"
)

type consumer interface {
	Use(middleware ...echo.MiddlewareFunc)
}

type namer interface {
	Name() string
}

var ErrUnableToExtractUserInfo = errors.New("unable to extract user info from context")

func ExtractUserInfo(c echo.Context) (*dtalk.UserTokenInfo, error) {
	info, ok := c.Get(UserInfoContextKey).(*dtalk.UserTokenInfo)
	if !ok {
		return nil, ErrUnableToExtractUserInfo
	}
	return info, nil
}

func logMiddlewareErr(c echo.Context, namer namer, err error) {
	log.Println(fmt.Errorf("%s middleware error on: %s: %w", namer, c.Path(), err))
}
