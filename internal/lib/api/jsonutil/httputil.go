package httputil

import (
	"drk-url-shortener/internal/lib/logger/sl"
	"encoding/json"
	"log/slog"
	"net/http"
)

// ReadJSON читает тело запроса и декодирует его в структуру.
// Также она автоматически закрывает тело запроса.
func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	defer r.Body.Close()

	// Ограничиваем размер тела (например, до 1 МБ) для защиты от DoS-атак
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	// Опционально: выдавать ошибку, если клиент прислал неизвестные поля
	dec.DisallowUnknownFields()

	return dec.Decode(data)
}

// WriteJSON сериализует данные в JSON и отправляет их клиенту с указанным HTTP-статусом.
func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	// Устанавливаем заголовок, что сервер возвращает JSON
	w.Header().Set("Content-Type", "application/json")

	// Устанавливаем HTTP-статус
	w.WriteHeader(status)

	// Записываем JSON напрямую в сетевой поток ответа
	return json.NewEncoder(w).Encode(data)
}

// HTTPError логирует ошибку через sl.Err и отправляет plain-text ответ клиенту.
func HTTPError(log *slog.Logger, w http.ResponseWriter, msg string, status int, err error) {
	log.Error(msg, sl.Err(err))
	http.Error(w, msg, status)
}
