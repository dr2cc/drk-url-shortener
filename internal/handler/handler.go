package handler

import (
	"drk-url-shortener/internal/config"
	"drk-url-shortener/internal/repository"

	"github.com/go-chi/chi/v5"
)

type server struct {
	router *chi.Mux
	cfg    config.Config
	repo   repository.Storage
}

func New(repo repository.Storage, cfg config.Config) *chi.Mux {
	// 3️⃣handler
	s := server{
		router: chi.NewRouter(),
		cfg:    cfg,
		repo:   repo,
	}

	// shortenText будет являтся методом server{}.
	// Метод shortenText сам будет является функцией http.HandlerFunc
	s.router.Post("/", s.shortenText)
	// Для примера redirect оставлю функцией, возвращающей хендлер (возвращает функцию http.HandlerFunc).
	// Функция выступают в роли фабрики (конструктора) для хендлера.
	// Она принимает специфичные параметры (s.repo)
	// и возвращает замыкание (переменной s.repo) в func(w http.ResponseWriter, r *http.Request) {...}
	s.router.Get("/{id}", redirect(s.repo))

	return s.router
}
