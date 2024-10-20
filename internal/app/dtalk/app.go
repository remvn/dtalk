package dtalk

type AppConfig struct {
	AppEnv               string `env:"APP_ENV"`
	AppPort              int    `env:"APP_PORT" envDefault:"8000"`
	JwtAccessTokenSecret string `env:"JWT_ACCESS_TOKEN_SECRET,required"`
	LiveKitHostURL       string `env:"LIVEKIT_HOST_URL,required"`
	LiveKitAPIKey        string `env:"LIVEKIT_API_KEY,required"`
	LiveKitAPISecret     string `env:"LIVEKIT_API_SECRET,required"`
	LiveKitClientURL     string `env:"LIVEKIT_CLIENT_URL,required"`
}
