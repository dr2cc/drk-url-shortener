package app

import (
	"drk-url-shortener/internal/config"
	"drk-url-shortener/internal/handler"
	"drk-url-shortener/internal/repository"
	"log"
	"net/http"
)

func Run(cfg config.Config) {
	// 1️⃣repository
	repo := repository.New()

	// 3️⃣handler
	mux := handler.New(repo, cfg)

	// Вторым параметром ListenAndServe получает:
	// mux (маршрутизатор= роутер= multiplexer) или
	// nil (используется маршрутизатор http.DefaultServeMux).
	// http.ListenAndServe(":8080", mux)

	err := http.ListenAndServe(cfg.ServAddres, mux)
	if err != nil {
		log.Fatalf("Start error: %s", err)
	}
}
