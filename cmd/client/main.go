package main

import (
	"drk-url-shortener/internal/cli"
	"log"
)

func main() {
	// Инициализируем клиент (парсим флаги, читаем конфиг)
	client := cli.New()

	// Запускаем клиентское приложение
	if err := client.Run(); err != nil {
		log.Fatalf("Client error: %v", err)
	}
}
