package handler

import (
	"drk-url-shortener/internal/config"

	"github.com/go-chi/chi/v5"
)

// hand - Snippet for http handler declaration

func New(repo map[string]string, cfg config.Config) *chi.Mux {
	// 3️⃣handler
	r := chi.NewRouter()
	r.Post("/", shortenText(repo, cfg))
	r.Get("/{id}", redirect(repo))
	return r
}
