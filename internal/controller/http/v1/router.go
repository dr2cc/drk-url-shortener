package v1

import (
	"drk-url-shortener/internal/config"
	"drk-url-shortener/internal/usecase"
	"log/slog"

	"github.com/go-chi/chi/v5"
	slogchi "github.com/samber/slog-chi"
)

// Stabilization Stage (Production-Ready MVP)
// 2. Переход на интерфейсы (Inversion of Control).

// Описываем структуру роутера для v1.
// Она инкапсулирует в себя зависимости, необходимые всем хэндлерам.
type Router struct {
	shortener *usecase.Shortener // Переходим на interface вместо конкретной структуры
	cfg       config.Config
	log       *slog.Logger
}

// NewRouter — конструктор, который настраивает маршруты для версии v1.
// Он возвращает готовый http.Handler, который можно подключить в главном app.go
func NewRouter(handler *chi.Mux, sh *usecase.Shortener, cfg config.Config, log *slog.Logger) *chi.Mux {
	// Инициализируем нашу внутреннюю структуру с зависимостями
	r := &Router{
		shortener: sh,
		cfg:       cfg,
		log:       log,
	}

	// Настраиваем middleware
	handler.Use(slogchi.New(log))

	// Привязываем эндпоинты напрямую к корню, как требует ТЗ
	handler.Post("/", r.shortenText)
	handler.Get("/{id}", r.redirect)

	// // Если нужно группировать маршруты, то выглядит примерно так:
	// handler.Route("/v1", func(chiRouter chi.Router) {
	// 	chiRouter.Post("/", r.shortenText)
	// 	chiRouter.Get("/{id}", r.redirect)
	// })
	return handler
}
