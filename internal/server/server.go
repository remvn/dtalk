package server

import (
	"fmt"

	"github.com/labstack/echo/v4"
	_ "github.com/livekit/protocol/livekit"
	_ "github.com/livekit/server-sdk-go/v2"
)

type Server struct {
	echoServer *echo.Echo
}

type LiveKitInfo struct {
	HostURL   string
	ApiKey    string
	ApiSecret string
}

func NewServer(info LiveKitInfo) *Server {
	echoServer := echo.New()
	server := &Server{
		echoServer: echoServer,
	}

	// hostURL := "ws://localhost:7880" // ex: https://project-123456.livekit.cloud
	// apiKey := "devkey"
	// apiSecret := "secret"
	// roomName := "myroom"
	// identity := "participantIdentity"
	// roomClient := lksdk.NewRoomServiceClient(hostURL, apiKey, apiSecret)

	return server
}

func (server *Server) Start(port int) {
	echoServer := server.echoServer
	echoServer.Logger.Fatal(echoServer.Start(fmt.Sprintf(":%d", port)))
}
