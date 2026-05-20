package handler

import (
	"drk-url-shortener/internal/config"
	"drk-url-shortener/internal/service"
	"log/slog"

	"github.com/go-chi/chi/v5"
	slogchi "github.com/samber/slog-chi"
)

// // URLSaverGetter описывает только то, что нужно хэндлерам
// type URLSaverGetter interface {
// 	SaveURL(alias string, url string) error
// 	GetURL(alias string) (string, error)
// }

type Handler struct {
	service *service.Service
}

// ❌Stabilization Stage (Production-Ready MVP)
// 2. Переход на интерфейсы (Inversion of Control).

func New(service *service.Service, cfg config.Config, log *slog.Logger) *chi.Mux {
	// 3️⃣handler
	r := chi.NewRouter()
	// Мы передаем настроенный logger внутрь middleware slog-chi
	r.Use(slogchi.New(log))
	r.Post("/", shortenText(service, cfg, log))
	r.Get("/{id}", redirect(service, log))
	return r
}
