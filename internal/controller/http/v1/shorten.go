package v1

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	httputil "drk-url-shortener/internal/lib/api/jsonutil"
)

// --- 1. СТРУКТУРЫ (DTO) ---

type ShortenRequest struct {
	URL string `json:"url"`
}

// Функция render.Bind(r *http.Request, v Binder) из пакета go-chi/render работает в два этапа:
// - Десериализация: автоматически определяет формат данных (JSON, XML) по заголовку Content-Type
// и декодирует тело запроса в переданную структуру.
// - Валидация: Если структура реализует метод Bind(r *http.Request) error,
// функция автоматически вызывает этот метод после декодирования.
// Если метод Bind возвращает ошибку, render.Bind прекращает обработку и возвращает эту ошибку в хендлер.
func (sr *ShortenRequest) Bind(r *http.Request) error {
	// 1. Очищаем пробелы (Sanitization)
	sr.URL = strings.TrimSpace(sr.URL)

	// 2. Проверяем, что URL вообще передан
	if sr.URL == "" {
		return errors.New("url field is required")
	}

	// 3. Проверяем, что это валидный URL-адрес
	u, err := url.ParseRequestURI(sr.URL)
	// строка URL обязательно должна содержать схему (протокол) и хост (доменное имя или IP-адрес)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return errors.New("invalid url format")
	}

	return nil // ошибок нет
}

type ShortenResponse struct {
	Result string `json:"result"`
}

// --- 2. ХЕНДЛЕРЫ ---

// // Пакетный JSON
// func (r Router) shortenBatch(w http.ResponseWriter, req *http.Request) {}

// Одиночный JSON
func (r Router) shortenJSON(w http.ResponseWriter, req *http.Request) {

	var sr ShortenRequest

	// 1. Декодируем JSON через стандартную библиотеку (тесты довольны!)
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields() // Опционально: запрещает лишние поля в JSON

	if err := decoder.Decode(&sr); err != nil {
		if errors.Is(err, io.EOF) {
			httputil.WriteJSONError(w, req, http.StatusBadRequest, "request body is empty") // пустое тело
			return
		}
		// "unexpected EOF" тут для совместимости с подходом на чистом Bind.
		// Но в целом логично: "понятный" EOF - пустое тело, остальное "unexpected EOF"
		httputil.WriteJSONError(w, req, http.StatusBadRequest, "unexpected EOF")
		return
	}

	// 2. "Вручную" вызываем метод валидации и очистки.
	if err := sr.Bind(req); err != nil {
		// Если JSON «битый» или не прошел валидацию в методе Bind (вернул ошибку)
		// err.Error() будет содержать то, что написано в unc (sr *ShortenRequest) Bind(r *http.Request) error{}
		// "url field is required" или "invalid url format"
		httputil.WriteJSONError(w, req, http.StatusBadRequest, err.Error())
		return
	}

	// // ВЕСЬ код выше (после объявления переменной можно удалить, но...
	// // "Чистый" переход не удался из-за автотестов Яндекса! Он не видит библиотеки дкодирования JSON.
	// // Переходим на Bind. render.Bind сам декодирует JSON и вызовет метод sr.Bind()
	// if err := render.Bind(req, &sr); err != nil {
	// 	// Если запрос пустой
	// 	if errors.Is(err, io.EOF) {
	// 		httputil.WriteJSONError(w, req, http.StatusBadRequest, "request body is empty")
	// 		return
	// 	}
	// 	// Если JSON «битый» или не прошел валидацию в методе Bind (вернул ошибку)
	// 	// err.Error() будет содержать то, что мы написали: "url field is required" или "invalid url format"
	// 	httputil.WriteJSONError(w, req, http.StatusBadRequest, err.Error())
	// 	return
	// }

	r.log.Info("request body decoded and validated", slog.Any("request", sr))

	// 3. Вызов бизнес-логики (UseCase)
	alias, err := r.shortener.Shorten(sr.URL)
	if err != nil {
		// Будущая проверка: если URL уже есть, возвращаем 409 Conflict
		// if errors.Is(err, domain.ErrURLExists) {
		//     if httputil.WriteJSONError(w, req, http.StatusConflict, "url already exists") { return }
		// }

		if httputil.WriteJSONError(w, req, http.StatusBadRequest, "failed to add url") {
			return
		}
	}

	r.log.Info("url added", slog.String("id", alias))

	// 4. Форматирование ответа
	response := ShortenResponse{
		Result: r.shortener.FormatShortURL(r.baseURL, alias),
	}

	// Отправляем успешный ответ клиенту. Чисто, надежно, в одну строку!
	httputil.JSON(w, req, http.StatusCreated, response)
}

// Текстовый хендлер.
// ♊Цепочка - 17.06.26 "GO десиарилизация. В проекте я использовал"
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
		if httputil.WriteTextError(r.log, w, req, "scroll reading error", http.StatusBadRequest, err) {
			return
		}
	}
	defer req.Body.Close() // Обязательно закрываем за посетителем дверь

	// 2. Проверяем содержимое (валидация)
	if len(body) == 0 {
		// Если свиток пуст — это Bad Request (ошибочный поиск)
		// Не можем применить HTTPError() - здесь не формируем ошибку!
		http.Error(w, "URL not found in request body", http.StatusBadRequest)
		return
	}

	// Stabilization Stage (Production-Ready MVP)
	// 1️⃣ Выделение слоев (Чистая архитектура)
	// 3. Обращение к UseCase/интерактору за алиасом (usecase и записывает его в db)
	alias, err := r.shortener.Shorten(string(body))
	if err != nil {
		if httputil.WriteTextError(r.log, w, req, "failed to add url", http.StatusBadRequest, err) {
			return
		}
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
