package main

import (
	"fmt"
	"net/http"
	"time"
)

// hand - Snippet for http handler declaration
func ShortenText(w http.ResponseWriter, r *http.Request) {
	// Здесь я должен получить текст из тела запроса (строку URL как text/plain)
	// и вернуть ответ с кодом 201 и сокращённым URL как text/plain
	// Видимо для начала разрешается фиксированная строка, к примеру "/EwHXdJfB"
	//
	// Все негативные кейсы- возвращаем 400
}

func Expand(w http.ResponseWriter, r *http.Request) {
	// Эндпоинт с методом GET и путём /{id}, где id — идентификатор сокращённого URL (например, /EwHXdJfB).
	// В случае успешной обработки запроса сервер возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
	//
	// Все негативные кейсы- возвращаем 400
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", ShortenText)
	mux.HandleFunc("GET /EwHXdJfB", Expand)

	// Вторым параметром ListenAndServe получает:
	// mux (маршрутизатор= роутер= multiplexer) или
	// nil (используется маршрутизатор http.DefaultServeMux).
	http.ListenAndServe(":8080", mux)
}
