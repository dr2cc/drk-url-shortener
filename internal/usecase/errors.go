package usecase

import "errors"

// Здесь ошибки usecase и repository
var (
	// "Чистая" ошибка бизнес-логики
	ErrCodeNotFound = errors.New("short code not found")
	// Техническая ошибка для репозитория
	ErrStorageNotFound = errors.New("entity not found in storage")
)
