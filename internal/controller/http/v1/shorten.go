package v1

import (
	"drk-url-shortener/internal/config"
	"drk-url-shortener/internal/service"
	"io"
	"log/slog"
	"net/http"
)

// Все негативные кейсы- возвращаем 400 = http.StatusBadRequest
func shortenText(service *service.Service, cfg config.Config, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Проверяем метод. В иудаике респонс — это ответ на вопрос.
		// Наш мудрец отвечает только на подношение данных (POST).
		if r.Method != http.MethodPost {
			http.Error(w, "The sage only accepts POST requests!", http.StatusBadRequest)
			return
		}

		// 2. Используем "слугу" io.LimitReader, чтобы подстраховаться.
		// Читаем не более 2 КБ, чтобы не переполнить "память".
		limitReader := io.LimitReader(r.Body, 2048)

		// Читаем из того, что "можно читать" (Reader)
		body, err := io.ReadAll(limitReader)
		if err != nil {
			http.Error(w, "Scroll reading error", http.StatusBadRequest)
			return
		}
		defer r.Body.Close() // Обязательно закрываем за посетителем дверь

		// 3. Проверяем содержимое (валидация)
		if len(body) == 0 {
			// Если свиток пуст — это Bad Request (ошибочный поиск)
			http.Error(w, "URL not found in request body", http.StatusBadRequest)
			return
		}

		// Stabilization Stage (Production-Ready MVP)
		// Выделение слоев (Чистая архитектура)
		alias, err := service.ShortenURL(string(body))
		if err != nil {
			log.Error("service.ShortenURL error:", "err", err)
			return
		}

		// 4. Формируем "Ответ-Обещание" (Response)
		// Сначала настраиваем "ящик" (ResponseWriter) в который будет положен respondēre- вердикт и ответ мудреца
		w.Header().Set("Content-Type", "text/plain")

		// Объявляем вердикт : "Создано" (201)
		w.WriteHeader(http.StatusCreated)

		// 4️⃣ Возвращаем клиенту response.
		// Готовим данные
		content := service.FormatShortURL(cfg.BaseURL, alias)

		// Из описания:
		// Функция Write записывает данные в соединение (to the connection) в HTTP-ответе.
		// Если метод ResponseWriter.WriteHeader еще не был вызван, Write вызывает WriteHeader(http.StatusOK) перед записью данных.
		// Если заголовок не содержит строку Content-Type (к примеру "w.Header().Set("Content-Type", "text/plain")"),
		// Write (при помощи DetectContentType) устанавливает Content-Type, по результату анализа начальных 512 байт возвращаемых данных.
		// Кроме того, если общий размер всех записанных данных составляет менее нескольких КБ и нет вызовов Flush,
		// заголовок Content-Length добавляется автоматически.
		//
		// Пишем (Write) в то, во что "можно писать" (...Writer)
		w.Write([]byte(content))
	}
}
