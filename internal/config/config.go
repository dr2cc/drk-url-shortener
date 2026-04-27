package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServAddres string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func New() *Config {
	cfg := &Config{}
	// Разбираем флаги в конфигурацию
	flag.StringVar(&cfg.ServAddres, "a", ":8080", "HTTP server startup address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base URL")
	flag.Parse()

	// // Удобная система, но не соответствует заданию yp.
	// // Загружаем переменные окружения из файла .env в корне
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatalf("Error loading env variables: %s", err.Error())
	// }

	// Разбираем переменные окружения в конфигурацию
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Config parsing error: %+v\n", err)
	}

	return cfg
}
