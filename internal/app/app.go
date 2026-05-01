package app

import (
	"drk-url-shortener/internal/config"
	"drk-url-shortener/internal/handler"
	"log"
	"net/http"
)

// 1️⃣repository
// var repo map[string]string

// var cfg config.Config

func Run(cfg config.Config) {
	// Будущая цепочка repository -> service -> handler

	// 1️⃣repository
	repo := make(map[string]string)

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
