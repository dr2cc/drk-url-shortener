package v1

import (
	"compress/flate"
	mw "drk-url-shortener/internal/controller/http/middleware"
	"drk-url-shortener/internal/usecase"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
)

// Stabilization Stage (Production-Ready MVP)
// 2️⃣ Осуществлен переход на интерфейсы (Inversion of Control).

// Описываем структуру роутера для v1.
// Она инкапсулирует в себя зависимости, необходимые всем хэндлерам.
type Router struct {
	shortener usecase.UseCase // interface вместо конкретной структуры
	baseURL   string
	log       *slog.Logger
}

// Принимай интерфейсы (usecase.UseCase), возвращай структуры (*chi.Mux)
func NewRouter(handler *chi.Mux, uc usecase.UseCase, baseURL string, log *slog.Logger) *chi.Mux {
	// Инициализируем нашу внутреннюю структуру с зависимостями
	r := &Router{
		shortener: uc,
		baseURL:   baseURL,
		log:       log,
	}

	// Настраиваем middleware
	handler.Use(slogchi.New(log))
	// Распаковываем входящий Gzip (если пришел gzip, он распаковывается и подменяет r.Body)
	handler.Use(mw.DecompressRequest)
	// Если клиент хочет сжатый ответ, запись в w перехватывается и сжимается (gzip)
	handler.Use(middleware.Compress(flate.BestSpeed))

	// Привязываем эндпоинты напрямую к корню (а не через v1), как требует ТЗ
	handler.Post("/", r.shortenText)
	handler.Get("/{id}", r.redirect)
	handler.Post("/api/shorten", r.shortenJSON)

	return handler
}
