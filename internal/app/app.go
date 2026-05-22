package app

import (
	"drk-url-shortener/internal/config"
	v1 "drk-url-shortener/internal/controller/http/v1"
	"drk-url-shortener/internal/lib/logger/sl"
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
	// (в yp это отдельный трек) Метрики: Интеграция с Prometheus для отслеживания количества запросов (RPS), времени ответа (latency) и количества ошибок 4xx / 5xx.
	log := sl.SetupLogger(cfg.Env)
	slog.SetDefault(log)
	log.Info("starting application", slog.String("env", cfg.Env))

	// 📍(сюда вернуться) Stabilization Stage (Production-Ready MVP)
	// 3️⃣ Инфраструктурный слой- db
	// Будет осуществлено в iter9 (сохранение сокращенных URL в файл при выходе и загрузка из него при запуске)) и затем в iter10 (pg)

	// handler
	r := chi.NewRouter()

	// DI
	repos := repository.New(make(map[string]string))
	// ⬇ Сервисам нужно то, что делает репозиторий (сохранение и нахождение).
	// 22.05.26 Продолжить от сюда https://share.google/aimode/ixhlMNK4DoCFiNBij
	services := usecase.New(repos)
	// ⬇ Хендлерам нужно то, что делает сервис (форматирование в результирующую строку, работа по сокращению, работа по получению).
	handlers := v1.NewRouter(r, services, cfg, log)

	// 📍(сюда вернуться) Stabilization Stage (Production-Ready MVP)
	// 5️⃣ Полноценная обработка контекста (context.Context)
	// Все сетевые запросы, походы в базу данных и логирование начнут использовать r.Context().
	// Это необходимо для graceful shutdown и для отмены долгих операций, если клиент разорвал соединение.
	log.Info("server is starting", "port", cfg.ServAddres)
	err := http.ListenAndServe(cfg.ServAddres, handlers)
	if err != nil {
		log.Error("start error:", "err", err)
		return err
	}
	return nil
}
