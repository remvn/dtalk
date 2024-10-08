package main

import (
	"dtalk/internal/config"
	"dtalk/internal/handler"
	"dtalk/internal/logic/lk"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type appConfig struct {
	Port              int    `env:"PORT" envDefault:"8080"`
	AccessTokenSecret string `env:"ACCESS_TOKEN_SECRET,required"`
	LiveKitHostURL    string `env:"LIVEKIT_HOST_URL,required"`
	LiveKitAPIKey     string `env:"LIVEKIT_API_KEY,required"`
	LiveKitAPISecret  string `env:"LIVEKIT_API_SECRET,required"`
}

func main() {
	_ = godotenv.Load()
	appConfig := new(appConfig)
	err := env.Parse(appConfig)
	if err != nil {
		log.Fatal(err)
	}

	info := lk.Config{
		HostURL:   appConfig.LiveKitHostURL,
		ApiKey:    appConfig.LiveKitAPIKey,
		ApiSecret: appConfig.LiveKitAPISecret,
	}

	config := handler.ServerConfig{
		AuthTokenConfig: config.JwtTokenConfig{
			Name:     "access_token",
			Secret:   []byte(appConfig.AccessTokenSecret),
			Duration: time.Hour * 24 * 10,
		},
	}
	server := handler.NewServer(config, info)

	server.Start(appConfig.Port)
}
