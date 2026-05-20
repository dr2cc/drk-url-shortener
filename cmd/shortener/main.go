package main

import (
	"drk-url-shortener/internal/app"
	"drk-url-shortener/internal/config"
	"fmt"
	"log"
	"os"
)

func main() {
	// Stabilization Stage (Production-Ready MVP)
	// 3️⃣Инфраструктурный слой- конфигурация.
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	if err := app.Run(cfg); err != nil {
		// Более сложная обработка ошибки
		// для дальнейшего внедрения graceful shutdown
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
