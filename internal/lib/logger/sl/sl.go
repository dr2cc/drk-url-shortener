package sl

import (
	"errors"
	"log/slog"
	"os"

	"github.com/Marlliton/slogpretty"
	"github.com/go-playground/validator/v10"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		// "Красивый" лог для разработки
		opts := slogpretty.DefaultOptions()
		log = slog.New(slogpretty.New(os.Stdout, opts))
	case envDev:
		// Обычный текстовый лог (если slogpretty не нужен на стейджинге (промежуточная тестовая среда,
		// максимально приближенная к продакшену))
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		// Строгий JSON для сервера
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

// // Простой хелпер
// func Err(err error) slog.Attr {
// 	return slog.Attr{
// 		Key:   "error",
// 		Value: slog.StringValue(err.Error()),
// 	}
// }

// Хелпер с учетом валидации json
func Err(err error) slog.Attr {
	var valErrs validator.ValidationErrors

	// Проверяем: если это ошибка валидации go-playground/validator
	if errors.As(err, &valErrs) {
		fields := make(map[string]string)
		for _, e := range valErrs {
			// Собираем карту: "URL": "required" или "URL": "url"
			fields[e.Field()] = e.Tag()
		}

		// Возвращаем структурированный объект (мапу) вместо плоской строки
		return slog.Attr{
			Key:   "validation_errors",
			Value: slog.AnyValue(fields),
		}
	}

	// Для всех остальных ошибок (БД, JSON, сеть) оставляем "стандартное" поведение
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
