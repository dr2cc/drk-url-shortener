package httputil

import (
	"drk-url-shortener/internal/lib/logger/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

type errorResponse struct {
	Message string `json:"message"`
}

// WriteJSONError теперь ЯВНО принимает логгер и НЕ возвращает bool
func WriteJSONError(log *slog.Logger, w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	// Использованием ErrorContext, чтобы логи связывались по Trace ID через контекст запроса
	log.ErrorContext(r.Context(), message,
		slog.Int("status", statusCode),
		slog.String("path", r.URL.Path),
	)

	render.Status(r, statusCode)
	render.JSON(w, r, errorResponse{Message: message})
}

// WriteTextError тоже использует ErrorContext для консистентности
func WriteTextError(log *slog.Logger, w http.ResponseWriter, r *http.Request, msg string, status int, err error) {
	// 1. Логируем ошибку через slog
	// Добавляем r.Context()
	log.ErrorContext(r.Context(), msg,
		sl.Err(err),
		slog.Int("status", status),
		slog.String("path", r.URL.Path), // полезно добавить и сюда
	)
	// // Прежний вариант
	// log.Error(msg, sl.Err(err))

	// 2. Оповещаем chi/render о статус-коде (важно для middleware-логгеров chi)
	render.Status(r, status)

	// 3. Отправляем plain-text через стандартный метод
	http.Error(w, msg, status)
}

// JSON — хелпер для отправки успешных ответов
func JSON(w http.ResponseWriter, r *http.Request, statusCode int, v interface{}) {
	render.Status(r, statusCode)
	render.JSON(w, r, v)
}
