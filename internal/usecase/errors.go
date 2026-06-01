package usecase

import "errors"

// Здесь ошибки usecase и repository
var (
	// Чистые бизнес-ошибки
	ErrCodeNotFound   = errors.New("short code not found")
	ErrInvalidLongURL = errors.New("provided URL is invalid")
	// Техническая ошибка для репозитория
	ErrStorageNotFound = errors.New("entity not found in storage")
)
