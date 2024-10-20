package main

import (
	"github.com/remvn/dtalk/internal/adapter/lk"
	"github.com/remvn/dtalk/internal/app/dtalk"
	"github.com/remvn/dtalk/internal/config"
	"github.com/remvn/dtalk/internal/setup"
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
			Secret:   []byte(appConfig.JwtAccessTokenSecret),
			Duration: time.Hour * 24 * 10,
		},
		CORS: cors,
	}
	server := setup.NewServer(serverConf, lkConf)

	server.Start(appConfig.AppPort)
}
