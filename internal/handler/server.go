package handler

import (
	"dtalk/internal/dtalk"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echoServer *echo.Echo
	LkService  *dtalk.LkService
}

func NewServer(lkOptions dtalk.LkOptions) *Server {
	echoServer := echo.New()
	lkService := dtalk.NewLkService(lkOptions)

	server := &Server{
		echoServer: echoServer,
		LkService:  lkService,
	}

	meetingHandler := NewMeetingHandler(echoServer, lkService)
	meetingHandler.Register()

	return server
}

func (server *Server) Start(port int) {
	echoServer := server.echoServer
	log.Fatal(echoServer.Start(fmt.Sprintf(":%d", port)))
}
