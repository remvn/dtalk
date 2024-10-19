package main

import (
	"dtalk/internal/adapter/lk"
	"dtalk/internal/app/dtalk"
	"dtalk/internal/config"
	"dtalk/internal/setup"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	appConfig := new(dtalk.AppConfig)
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

	// setup server a->z
	serverConf := setup.ServerConfig{
		AppConfig: *appConfig,
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
