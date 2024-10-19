package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type PublicHandler struct {
	echoServer       *echo.Echo
	livekitClientURL string
}

func NewPublicHandler(
	echoServer *echo.Echo,
	livekitClientURL string,
) *PublicHandler {
	handler := &PublicHandler{
		echoServer:       echoServer,
		livekitClientURL: livekitClientURL,
	}
	return handler
}

func (handler *PublicHandler) Register(parentGroup *echo.Group) {
	group := parentGroup.Group("/public")
	group.GET("/livekit-client-url", handler.LivekitClientURL)
}

type liveKitClientURLRes struct {
	URL string `json:"url"`
}

func (handler *PublicHandler) LivekitClientURL(c echo.Context) error {
	return c.JSON(http.StatusOK, liveKitClientURLRes{
		URL: handler.livekitClientURL,
	})
}
