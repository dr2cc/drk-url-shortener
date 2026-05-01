package main

import (
	"drk-url-shortener/internal/app"
	"drk-url-shortener/internal/config"
)

func main() {
	cfg := config.New()

	app.Run(cfg)
}
