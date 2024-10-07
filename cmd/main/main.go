package main

import (
	"dtalk/internal/dtalk"
	"dtalk/internal/handler"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("unable to get PORT env")
	}
	_ = os.Getenv("APP_ENV")

	info := dtalk.LkOptions{
		HostURL:   os.Getenv("LIVEKIT_HOST_URL"),
		ApiKey:    os.Getenv("LIVEKIT_API_KEY"),
		ApiSecret: os.Getenv("LIVEKIT_API_SECRET"),
	}
	server := handler.NewServer(info)

	server.Start(port)
}
