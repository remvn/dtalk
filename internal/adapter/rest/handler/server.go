package handler

import (
	"dtalk/internal/adapter/rest/middleware"
	"dtalk/internal/app/logic/lk"
	"dtalk/internal/config"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
)

type Server struct {
	echoServer *echo.Echo
	LkService  *lk.Service
}

type ServerConfig struct {
	AuthTokenConfig config.JwtTokenConfig
	CORS            []string
}

func NewServer(config ServerConfig, lkConfig lk.Config) *Server {
	echoServer := echo.New()
	if len(config.CORS) > 0 {
		echoServer.Use(echoMW.CORSWithConfig(echoMW.CORSConfig{
			AllowOrigins: config.CORS,
		}))
	}

	lkService := lk.NewLkService(lkConfig)

	server := &Server{
		echoServer: echoServer,
		LkService:  lkService,
	}

	parentGroup := echoServer.Group("/api")
	authMiddleware := middleware.NewAuthMiddleware(config.AuthTokenConfig)

	authHandler := NewAuthHandler(
		config.AuthTokenConfig,
		echoServer,
		authMiddleware,
	)
	authHandler.Register(parentGroup)

	meetingHandler := NewMeetingHandler(
		echoServer,
		lkService,
		authMiddleware,
	)
	meetingHandler.Register(parentGroup)

	return server
}

func (server *Server) Start(port int) {
	echoServer := server.echoServer
	log.Fatal(echoServer.Start(fmt.Sprintf(":%d", port)))
}
