package setup

import (
	"dtalk/internal/adapter/lk"
	"dtalk/internal/adapter/rest/handler"
	"dtalk/internal/adapter/rest/middleware"
	"dtalk/internal/app/logic/meeting"
	"dtalk/internal/app/port"
	"dtalk/internal/config"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	echoMW "github.com/labstack/echo/v4/middleware"
)

// wire up everything

type Server struct {
	echoServer *echo.Echo

	roomManager    port.RoomManager
	meetingService port.MeetingPort
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

	roomManager := lk.NewLkRoomManager(lkConfig)
	meetingService := meeting.NewMeetingService(roomManager)

	server := &Server{
		echoServer: echoServer,

		roomManager:    roomManager,
		meetingService: meetingService,
	}

	parentGroup := echoServer.Group("/api")
	authMiddleware := middleware.NewAuth(config.AuthTokenConfig)
	roomAuthMiddleware := middleware.NewRoomAuth(meetingService)

	authHandler := handler.NewAuthHandler(
		config.AuthTokenConfig,
		echoServer,
		authMiddleware,
	)
	authHandler.Register(parentGroup)

	meetingHandler := handler.NewMeetingHandler(
		echoServer,
		authMiddleware,
		roomAuthMiddleware,
		meetingService,
	)
	meetingHandler.Register(parentGroup)

	return server
}

func (server *Server) Start(port int) {
	echoServer := server.echoServer
	log.Fatal(echoServer.Start(fmt.Sprintf(":%d", port)))
}
