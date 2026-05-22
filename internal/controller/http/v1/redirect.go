package v1

import (
	"drk-url-shortener/internal/service"
	"log/slog"
	"net/http"
)

// Все негативные кейсы- возвращаем 400 = http.StatusBadRequest
func redirect(service *service.Service, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// // Ниже- родной для chi метод определения id
		// // Но с ним не работают простые (и универсальные) тесты
		// id := chi.URLParam(r, "id")
		// Стандартный для встроенного роутера, должен поддерживаться chi в 2026
		id := r.PathValue("id")

		if r.Method != http.MethodGet {
			http.Error(w, "accepts GET requests!", http.StatusBadRequest)
			return
		}

		// Здесь обращение к сервису за url

		url, err := service.FindURL(id)
		if err != nil {
			log.Error("service.FindURL error:", "err", err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		// fmt.Fprintf(w, "www.google.com %s", time.Now())
	}
}
