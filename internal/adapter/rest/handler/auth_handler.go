package handler

import (
	"dtalk/internal/adapter/rest/middleware"
	"dtalk/internal/config"
	"dtalk/internal/pkg/random"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	tokenConfig    config.JwtTokenConfig
	echoServer     *echo.Echo
	authMiddleware *middleware.AuthMiddleware
}

func NewAuthHandler(
	tokenConfig config.JwtTokenConfig,
	echoServer *echo.Echo,
	authMiddleware *middleware.AuthMiddleware,
) *AuthHandler {
	handler := &AuthHandler{
		echoServer:     echoServer,
		tokenConfig:    tokenConfig,
		authMiddleware: authMiddleware,
	}
	return handler
}

func (handler *AuthHandler) Register(parentGroup *echo.Group) {
	group := parentGroup.Group("/auth")
	group.POST("/request-token", handler.requestToken)
}

type requestTokenDto struct {
	Name string `json:"name"`
}

type requestTokenRes struct {
	AccessToken string `json:"access_token"`
}

func (handler *AuthHandler) requestToken(c echo.Context) error {
	dto := &requestTokenDto{}
	if err := c.Bind(dto); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	claims := jwt.MapClaims{
		"name": dto.Name,
		"id":   random.GenerateID(),
	}
	tokenStr, err := handler.tokenConfig.Sign(claims)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// c.SetCookie(&http.Cookie{
	// 	Name:     handler.tokenConfig.Name,
	// 	Value:    tokenStr,
	// 	HttpOnly: true,
	// 	Expires:  handler.tokenConfig.GetExpire(),
	// })

	return c.JSON(http.StatusOK, requestTokenRes{
		AccessToken: tokenStr,
	})
}
