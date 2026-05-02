package handler

import (
	"drk-url-shortener/internal/config"
	"log/slog"

	"github.com/go-chi/chi/v5"
	slogchi "github.com/samber/slog-chi"
)

// hand - Snippet for http handler declaration

func New(repo map[string]string, cfg config.Config, log *slog.Logger) *chi.Mux {
	// 3️⃣handler
	r := chi.NewRouter()
	// Мы передаем настроенный logger внутрь middleware slog-chi
	r.Use(slogchi.New(log))
	r.Post("/", shortenText(repo, cfg))
	r.Get("/{id}", redirect(repo))
	return r
}
