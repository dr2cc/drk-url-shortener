package main

import (
	"drk-url-shortener/internal/config"
	"drk-url-shortener/internal/lib/random"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
)

// hand - Snippet for http handler declaration

// TODO: move to config
const aliasLength = 6

// 1️⃣repository
var repo map[string]string

var cfg config.Config

// Все негативные кейсы- возвращаем 400 = http.StatusBadRequest
func ShortenText(w http.ResponseWriter, r *http.Request) {
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

	// 2️⃣service
	alias := random.NewRandomString(aliasLength)
	// Запись в db
	repo[alias] = string(body)

	// 4. Формируем "Ответ-Обещание" (Response)
	// Сначала настраиваем "ящик" (ResponseWriter) в который будет положен respondēre- вердикт и ответ мудреца
	w.Header().Set("Content-Type", "text/plain")

	// Объявляем вердикт : "Создано" (201)
	w.WriteHeader(http.StatusCreated)

	// Из описаеия:
	// Функция Write записывает данные в соединение (to the connection) в HTTP-ответе.
	// Если метод ResponseWriter.WriteHeader еще не был вызван, Write вызывает WriteHeader(http.StatusOK) перед записью данных.
	// Если заголовок не содержит строку Content-Type (к примеру "w.Header().Set("Content-Type", "text/plain")"),
	// Write (при помощи DetectContentType) устанавливает Content-Type, по результату анализа начальных 512 байт возвращаемых данных.
	// Кроме того, если общий размер всех записанных данных составляет менее нескольких КБ и нет вызовов Flush,
	// заголовок Content-Length добавляется автоматически.
	//
	// Пишем (Write) в то, во что "можно писать" (...Writer)
	w.Write([]byte(cfg.BaseURL + "/" + alias))
}

// Все негативные кейсы- возвращаем 400 = http.StatusBadRequest
func Expand(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, repo[id], http.StatusTemporaryRedirect)
	// fmt.Fprintf(w, "www.google.com %s", time.Now())
}

func main() {
	flag.StringVar(&cfg.ServAddres, "a", ":8080", "HTTP server startup address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base URL")
	flag.Parse()

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Config parsing error: %+v\n", err)
	}

	// Будущая цепочка repository -> service -> handler
	repo = make(map[string]string)

	// 3️⃣handler
	mux := chi.NewRouter()
	mux.Post("/", ShortenText)
	mux.Get("/{id}", Expand)

	// Вторым параметром ListenAndServe получает:
	// mux (маршрутизатор= роутер= multiplexer) или
	// nil (используется маршрутизатор http.DefaultServeMux).
	// http.ListenAndServe(":8080", mux)

	err := http.ListenAndServe(cfg.ServAddres, mux)
	if err != nil {
		log.Fatalf("Start error: %s", err)
	}
}
