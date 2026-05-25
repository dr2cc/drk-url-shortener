package v1

import (
	"net/http"
)

func (r router) redirect(w http.ResponseWriter, req *http.Request) {
	// // Ниже- родной для chi метод определения id
	// // Но с ним не работают простые (и универсальные) тесты
	// id := chi.URLParam(r, "id")
	// Стандартный для встроенного роутера, должен поддерживаться chi в 2026
	id := req.PathValue("id")

	if req.Method != http.MethodGet {
		http.Error(w, "accepts GET requests!", http.StatusBadRequest)
		return
	}

	// Здесь обращение к сервису за url

	url, err := r.shortener.GetOriginal(id)
	if err != nil {
		r.log.Error("service.FindURL error:", "err", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
	// fmt.Fprintf(w, "www.google.com %s", time.Now())
}

// // Все негативные кейсы- возвращаем 400 = http.StatusBadRequest
// func redirect(service *service.Service, log *slog.Logger) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// // Ниже- родной для chi метод определения id
// 		// // Но с ним не работают простые (и универсальные) тесты
// 		// id := chi.URLParam(r, "id")
// 		// Стандартный для встроенного роутера, должен поддерживаться chi в 2026
// 		id := r.PathValue("id")

// 		if r.Method != http.MethodGet {
// 			http.Error(w, "accepts GET requests!", http.StatusBadRequest)
// 			return
// 		}

// 		// Здесь обращение к сервису за url

// 		url, err := service.FindURL(id)
// 		if err != nil {
// 			log.Error("service.FindURL error:", "err", err)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "text/plain")
// 		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
// 		// fmt.Fprintf(w, "www.google.com %s", time.Now())
// 	}
// }
