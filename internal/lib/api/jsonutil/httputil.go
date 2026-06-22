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

type statusResponse struct {
	Status string `json:"status"`
}

// WriteJSONError логирует ошибку через slog, отправляет JSON и возвращает true
func WriteJSONError(w http.ResponseWriter, r *http.Request, statusCode int, message string) bool {
	// Использование структурированного логирования slog
	slog.Error(message,
		slog.Int("status", statusCode),
		slog.String("path", r.URL.Path),
	)

	render.Status(r, statusCode)
	render.JSON(w, r, errorResponse{Message: message})

	return true
}

// WriteTextError логирует ошибку через slog, отправляет plain-text ответ клиенту и возвращает true.
func WriteTextError(log *slog.Logger, w http.ResponseWriter, r *http.Request, msg string, status int, err error) bool {
	// 1. Логируем ошибку через slog
	log.Error(msg, sl.Err(err))

	// 2. Оповещаем chi/render о статус-коде (важно для middleware-логгеров chi)
	render.Status(r, status)

	// 3. Отправляем plain-text через стандартный метод
	http.Error(w, msg, status)

	return true
}

// JSON — хелпер для отправки успешных ответов
func JSON(w http.ResponseWriter, r *http.Request, statusCode int, v interface{}) {
	render.Status(r, statusCode)
	render.JSON(w, r, v)
}
