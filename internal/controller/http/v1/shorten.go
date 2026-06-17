package v1

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	resp "drk-url-shortener/internal/lib/api/response"
	"drk-url-shortener/internal/lib/logger/sl"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type ShortenRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

func (r Router) shortenJSON(w http.ResponseWriter, req *http.Request) {

	var sr ShortenRequest

	err := render.DecodeJSON(req.Body, &sr)
	if errors.Is(err, io.EOF) {
		// Такую ошибку встретим, если получили запрос с пустым телом.
		// Обработаем её отдельно
		r.log.Error("request body is empty")

		http.Error(w, "URL not found in request body", http.StatusBadRequest)
		render.JSON(w, req, resp.Error("empty request"))

		return
	}
	if err != nil {
		r.log.Error("failed to decode request body", sl.Err(err))

		render.JSON(w, req, resp.Error("failed to decode request"))

		return
	}

	r.log.Info("request body decoded", slog.Any("request", sr))
	if err := validator.New().Struct(sr); err != nil {
		validateErr := err.(validator.ValidationErrors)

		r.log.Error("invalid request", sl.Err(err))

		render.JSON(w, req, resp.ValidationError(validateErr))

		return
	}

	alias, err := r.shortener.Shorten(sr.URL)
	if err != nil {
		// Полная ошибка в лог
		r.log.Error("failed to add url", sl.Err(err))

		render.JSON(w, req, resp.Error("failed to add url"))
		// http.Error(w, "failed to add url", http.StatusBadRequest)
		return
	}
	r.log.Info("url added", slog.String("id", alias))

	// responseOK(w, req, alias)

	// }

	// func responseOK(w http.ResponseWriter, r *http.Request, alias string) {

	// 📌render.JSON устанавливает нужный Content-Type, но нам нужен http.StatusCreated
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Обращение к UseCase/интерактору за форматированем
	content := r.shortener.FormatShortURL(r.baseURL, alias)
	render.JSON(w, req, ShortenResponse{
		Result: content,
	})
}

// Все негативные кейсы- возвращаем 400 = http.StatusBadRequest
func (r Router) shortenText(w http.ResponseWriter, req *http.Request) {
	// // Эта проверка внутри хендлера вредна, и её нужно удалить по двум причинам:
	// // - Нарушение ответственности (SRP): Фильтрация методов — это задача роутера, а не бизнес-логики хендлера.
	// // - Мертвый код (Dead Code): Из-за того, что роутер handler.Post уже фильтрует трафик,
	// // условие if req.Method != http.MethodPost никогда не выполнится.
	// if req.Method != http.MethodPost {
	// 	http.Error(w, "The sage only accepts POST requests!", http.StatusBadRequest)
	// 	return
	// }

	// // Используем "слугу" io.LimitReader, чтобы подстраховаться.
	// // СМЫСЛ понимаю, но не реализацию (особенно в тесте). Пока не делаю..
	// limitReader := io.LimitReader(req.Body, 2048)

	// 1. Читаем (тело) из того, что "можно читать" (Reader)
	body, err := io.ReadAll(req.Body) // (limitReader)
	if err != nil {
		// scroll- свиток!
		http.Error(w, "scroll reading error", http.StatusBadRequest)
		return
	}
	defer req.Body.Close() // Обязательно закрываем за посетителем дверь

	// 2. Проверяем содержимое (валидация)
	if len(body) == 0 {
		// Если свиток пуст — это Bad Request (ошибочный поиск)
		http.Error(w, "URL not found in request body", http.StatusBadRequest)
		return
	}

	// Stabilization Stage (Production-Ready MVP)
	// 1️⃣ Выделение слоев (Чистая архитектура)
	// 3. Обращение к UseCase/интерактору за алиасом (usecase и записывает его в db)
	alias, err := r.shortener.Shorten(string(body))
	if err != nil {
		// Полная ошибка в лог
		r.log.Error("failed to add url", "err", err)
		http.Error(w, "failed to add url", http.StatusBadRequest)
		return
	}

	// 4. Формируем "Ответ-Обещание" (Response)
	// Сначала настраиваем "ящик" (ResponseWriter) в который будет положен respondēre- вердикт и ответ мудреца
	w.Header().Set("Content-Type", "text/plain")

	// Объявляем вердикт : "Создано" (201)
	w.WriteHeader(http.StatusCreated)

	// 5. Возвращаем клиенту response.
	// Готовим данные
	// Обращение к UseCase/интерактору за форматированем
	content := r.shortener.FormatShortURL(r.baseURL, alias)

	// Пишем (Write) в то, во что "можно писать" (...Writer)
	w.Write([]byte(content))
	// Из описания:
	// Функция Write записывает данные в соединение (to the connection) в HTTP-ответе.
	// Если метод ResponseWriter.WriteHeader еще не был вызван, Write вызывает WriteHeader(http.StatusOK) перед записью данных.
	// Если заголовок не содержит строку Content-Type (к примеру "w.Header().Set("Content-Type", "text/plain")"),
	// Write (при помощи DetectContentType) устанавливает Content-Type, по результату анализа начальных 512 байт возвращаемых данных.
	// Кроме того, если общий размер всех записанных данных составляет менее нескольких КБ и нет вызовов Flush,
	// заголовок Content-Length добавляется автоматически.
}
