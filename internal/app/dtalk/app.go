package dtalk

type AppConfig struct {
	AppEnv            string `env:"APP_ENV"`
	Port              int    `env:"PORT" envDefault:"8080"`
	AccessTokenSecret string `env:"ACCESS_TOKEN_SECRET,required"`
	LiveKitHostURL    string `env:"LIVEKIT_HOST_URL,required"`
	LiveKitAPIKey     string `env:"LIVEKIT_API_KEY,required"`
	LiveKitAPISecret  string `env:"LIVEKIT_API_SECRET,required"`
	LiveKitClientURL  string `env:"LIVEKIT_CLIENT_URL,required"`
}
