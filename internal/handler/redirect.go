package handler

import (
	"drk-url-shortener/internal/repository"
	"net/http"
)

// Все негативные кейсы- возвращаем 400 = http.StatusBadRequest
func redirect(repo repository.Storage) http.HandlerFunc {
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

		w.Header().Set("Content-Type", "text/plain")

		slug, err := repo.Read(id)
		// Сейчас не срабатывает (т.к. repo.Read всегда возвращает nil)
		// Но yp тест должен пройти!
		if err != nil {
			http.Error(w, "repo.Read error", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, slug, http.StatusTemporaryRedirect)
		// fmt.Fprintf(w, "www.google.com %s", time.Now())
	}
}
