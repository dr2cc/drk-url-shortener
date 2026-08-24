package app

import (
	"drk-url-shortener/internal/config"
	v1 "drk-url-shortener/internal/controller/http/v1"
	"drk-url-shortener/internal/lib/logger/sl"
	"drk-url-shortener/internal/lib/random"
	"drk-url-shortener/internal/repository"
	"drk-url-shortener/internal/usecase"

	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run(cfg config.Config) error {
	// 📌Stabilization Stage (Production-Ready MVP)
	// 4️⃣ Наблюдаемость (Observability).
	// Логер добавден.
	// Еще следует добавить:
	// Трейсинг: Внедрение OpenTelemetry (Jaeger) для отслеживания пути запроса, особенно когда проект начнет ходить в базу данных.
	// Метрики (в yp это отдельный трек): Интеграция с Prometheus для отслеживания количества запросов (RPS), времени ответа (latency) и количества ошибок 4xx / 5xx.
	logger := sl.SetupLogger(cfg.Env)
	// Главное назначение slog.SetDefault(logger) это переключение на работу в режиме нового slog во всем приложении (даже при использовании log.).
	// Без него, если использовать slog. или log. то будет выводиться текст (старый log).
	// Или придется всюду передавать переменную logger
	slog.SetDefault(logger)
	slog.Info("starting application", slog.String("env", cfg.Env))

	// 📍(сюда вернуться) Stabilization Stage (Production-Ready MVP)
	// 3️⃣ Инфраструктурный слой- db
	// Будет осуществлено в iter9 (сохранение сокращенных URL в файл при выходе и загрузка из него при запуске))
	// и затем в iter10 (pg)

	// mux
	mux := chi.NewRouter()

	// DI
	generator := random.NewBase62Generator() // генератор кодов

	// Выбор между именованным полем структуры и неименованным (встраиванием/embedding)  https://share.google/aimode/ixhlMNK4DoCFiNBij
	// Using the Factory Pattern
	repos, err := repository.New(cfg.CacheDumpPath)
	if err != nil {
		// log.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
		slog.Error("app - Run - repository.New:", "err", err)
		return err
	}
	// ⬇ Сервисам нужно то, что делает репозиторий (сохранение и нахождение).
	shortenerUseCase := usecase.New(repos, generator)
	// ⬇ Хендлерам нужно то, что делает сервис (форматирование, работа по сокращению, работа по получению).
	handlers := v1.NewRouter(mux, shortenerUseCase, cfg.BaseURL, logger)

	// 📍(сюда вернуться) Stabilization Stage (Production-Ready MVP)
	// 5️⃣ Полноценная обработка контекста. Должны начать использовать контекст:
	// 1. Сетевые запросы;
	// 2. "Походы" в базу данных.
	// Это необходимо для graceful shutdown и для отмены долгих операций, если клиент разорвал соединение.
	slog.Info("server is starting", "port", cfg.ServAddres)
	err = http.ListenAndServe(cfg.ServAddres, handlers)
	if err != nil {
		slog.Error("start error:", "err", err)
		return err
	}
	return nil
}
