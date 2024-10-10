package handler

import (
	"dtalk/internal/config"
	"dtalk/internal/logic/lk"
	"dtalk/internal/middleware"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echoServer *echo.Echo
	LkService  *lk.Service
}

type ServerConfig struct {
	AuthTokenConfig config.JwtTokenConfig
}

func NewServer(config ServerConfig, lkConfig lk.Config) *Server {
	echoServer := echo.New()
	lkService := lk.NewLkService(lkConfig)

	server := &Server{
		echoServer: echoServer,
		LkService:  lkService,
	}

	authMiddleware := middleware.NewAuthMiddleware(config.AuthTokenConfig)

	authHandler := NewAuthHandler(
		config.AuthTokenConfig,
		echoServer,
		authMiddleware,
	)
	authHandler.Register()

	meetingHandler := NewMeetingHandler(
		echoServer,
		lkService,
		authMiddleware,
	)
	meetingHandler.Register()

	return server
}

func (server *Server) Start(port int) {
	echoServer := server.echoServer
	log.Fatal(echoServer.Start(fmt.Sprintf(":%d", port)))
}
