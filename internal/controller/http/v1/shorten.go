package v1

import (
	"io"
	"net/http"
)

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
