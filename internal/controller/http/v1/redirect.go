package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (r Router) redirect(w http.ResponseWriter, req *http.Request) {
	// Ниже- родной для chi метод определения id
	// Но с ним не работают простые (и универсальные) тесты
	id := chi.URLParam(req, "id")
	// // Стандартный для встроенного роутера, должен поддерживаться chi в 2026
	// id := req.PathValue("id")

	if req.Method != http.MethodGet {
		http.Error(w, "accepts GET requests!", http.StatusBadRequest)
		return
	}

	// Обращение к сервису за url
	url, err := r.shortener.GetOriginal(id)

	if err != nil {
		r.log.Error("shortener.GetOriginal error:", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
	// fmt.Fprintf(w, "www.google.com %s", time.Now())
}
