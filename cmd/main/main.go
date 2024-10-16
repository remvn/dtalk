package main

import (
	"dtalk/internal/adapter/lk"
	"dtalk/internal/config"
	"dtalk/internal/setup"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type appConfig struct {
	AppEnv            string `env:"APP_ENV"`
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

	if appConfig.AppEnv != "local" && appConfig.AppEnv != "production" {
		log.Fatalf("invalid APP_ENV: %s", appConfig.AppEnv)
	}

	lkConf := lk.Config{
		HostURL:   appConfig.LiveKitHostURL,
		ApiKey:    appConfig.LiveKitAPIKey,
		ApiSecret: appConfig.LiveKitAPISecret,
	}

	cors := []string{}
	if appConfig.AppEnv == "local" {
		cors = []string{"*"}
	}

	serverConf := setup.ServerConfig{
		AuthTokenConfig: config.JwtTokenConfig{
			Name:     "access_token",
			Secret:   []byte(appConfig.AccessTokenSecret),
			Duration: time.Hour * 24 * 10,
		},
		CORS: cors,
	}
	server := setup.NewServer(serverConf, lkConf)

	server.Start(appConfig.Port)
}
