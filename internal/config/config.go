package config

type Config struct {
	ServAddres string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}
