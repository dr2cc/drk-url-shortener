package app

import (
	"drk-url-shortener/internal/config"
	"drk-url-shortener/internal/handler"
	"drk-url-shortener/internal/lib/logger/sl"
	"log/slog"
	"net/http"
	"os"
)

func Run(cfg config.Config) {
	log := sl.SetupLogger(cfg.Env)
	slog.SetDefault(log)
	log.Info("starting application", slog.String("env", cfg.Env))

	// Будущая цепочка repository -> service -> handler

	// 1️⃣repository
	repo := make(map[string]string)

	// 3️⃣handler
	mux := handler.New(repo, cfg, log)

	// Вторым параметром ListenAndServe получает:
	// mux (маршрутизатор= роутер= multiplexer) или
	// nil (используется маршрутизатор http.DefaultServeMux).
	// http.ListenAndServe(":8080", mux)
	log.Info("server is starting", "port", cfg.ServAddres)
	err := http.ListenAndServe(cfg.ServAddres, mux)
	if err != nil {
		// log.Fatalf("Start error: %s", err)
		// slog аналог для log.Fatal(err)
		log.Error("start error:", "err", err)
		os.Exit(1)
	}
}
