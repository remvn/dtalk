package middleware

import (
	"dtalk/internal/dtalk"
	"errors"

	"github.com/labstack/echo/v4"
)

const (
	UserTokenContextKey = "user_claims"
	UserInfoContextKey  = "user_info"
)

type consumer interface {
	Use(middleware ...echo.MiddlewareFunc)
}

var ErrUnableToExtractUserInfo = errors.New("unable to extract user info from context")

func ExtractUserInfo(c echo.Context) (*dtalk.UserTokenInfo, error) {
	info, ok := c.Get(UserInfoContextKey).(*dtalk.UserTokenInfo)
	if !ok {
		return nil, ErrUnableToExtractUserInfo
	}
	return info, nil
}
