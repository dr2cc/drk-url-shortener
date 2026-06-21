package v1

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	httputil "drk-url-shortener/internal/lib/api/jsonutil"
	"drk-url-shortener/internal/lib/logger/sl"
)

// --- 1. СТРУКТУРЫ (DTO) ---

type ShortenRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type ShortenResponse struct {
	Result string `json:"result"`
}

// --- 2. ХЕНДЛЕРЫ ---

// // Пакетный JSON
// func (r Router) shortenBatchHandler(w http.ResponseWriter, req *http.Request) {}

// Одиночный JSON
func (r Router) shortenJSON(w http.ResponseWriter, req *http.Request) {

	var sr ShortenRequest
	// Собственная функция добавляющая некоторые проверки и обработку ошибок.
	err := httputil.ReadJSON(w, req, &sr)
	if errors.Is(err, io.EOF) {
		// Обработаем отдельно ошибку, если получили запрос с пустым телом.
		httputil.HTTPError(r.log, w, "request body is empty", http.StatusBadRequest, err)
		return
	}
	if err != nil {
		httputil.HTTPError(r.log, w, "failed to decode request body", http.StatusBadRequest, err)
		return
	}

	r.log.Info("request body decoded", slog.Any("request", sr))
	if err := r.validator.Struct(sr); err != nil {
		//if err := validator.New().Struct(sr); err != nil {
		// // Пишут, что создавать валидатор через validator.New()
		// // прямо внутри хендлера на каждый запрос — это плохая практика (антипаттерн).
		// // Метод validator.New() под капотом делает очень много «тяжелой» работы!
		httputil.HTTPError(r.log, w, "invalid request", http.StatusBadRequest, err)
		return
	}

	// Обращение к UseCase/интерактору за alias
	alias, err := r.shortener.Shorten(sr.URL)
	if err != nil {
		httputil.HTTPError(r.log, w, "failed to add url", http.StatusBadRequest, err)
		return
	}

	r.log.Info("url added", slog.String("id", alias))

	// Обращение к UseCase/интерактору за форматированем
	response := ShortenResponse{
		Result: r.shortener.FormatShortURL(r.baseURL, alias),
	}

	// 2. Записываем ответ со статусом 201 (Created)
	if err := httputil.WriteJSON(w, http.StatusCreated, response); err != nil {
		// Если не удалось записать ответ (например, клиент отключился)
		// Здесь только лог, так как заголовки уже отправлены
		r.log.Error("failed to write response", sl.Err(err))
	}
}

// Текстовый хендлер
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
		httputil.HTTPError(r.log, w, "scroll reading error", http.StatusBadRequest, err)
		return
	}
	defer req.Body.Close() // Обязательно закрываем за посетителем дверь

	// 2. Проверяем содержимое (валидация)
	if len(body) == 0 {
		// Если свиток пуст — это Bad Request (ошибочный поиск)
		// Не можем применить httputil.HTTPError() - здесь не формируем ошибку!
		http.Error(w, "URL not found in request body", http.StatusBadRequest)
		return
	}

	// Stabilization Stage (Production-Ready MVP)
	// 1️⃣ Выделение слоев (Чистая архитектура)
	// 3. Обращение к UseCase/интерактору за алиасом (usecase и записывает его в db)
	alias, err := r.shortener.Shorten(string(body))
	if err != nil {
		httputil.HTTPError(r.log, w, "failed to add url", http.StatusBadRequest, err)
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
