package v1

import (
	"drk-url-shortener/internal/usecase"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (r Router) redirect(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	if req.Method != http.MethodGet {
		http.Error(w, "accepts GET requests!", http.StatusBadRequest)
		return
	}

	// Обращение к UseCase/интерактору за url
	url, err := r.shortener.GetOriginal(id)
	if err != nil {
		// 1. Всегда логируем полную техническую ошибку
		r.log.Error("shortener.GetOriginal error", "err", err, "id", id)

		// 2. Проверяем бизнес-ошибку для пользователя
		if errors.Is(err, usecase.ErrCodeNotFound) {
			// Возвращаем понятный текст и правильный статус-код по ТЗ (400)
			http.Error(w, "short link not found", http.StatusBadRequest)
			return
		}

		// 3. Для всех остальных неизвестных ошибок (упала база, сеть и т.д.)
		// должны отдавать стандартный http.StatusInternalServerError, но по ТЗ (400)
		http.Error(w, "internal server error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}
