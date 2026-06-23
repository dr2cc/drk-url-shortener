package sl

import (
	"log/slog"
	"os"

	"github.com/Marlliton/slogpretty"
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

// Простой хелпер, для shortenText
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
