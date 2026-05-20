package app

import (
	"drk-url-shortener/internal/config"
	"drk-url-shortener/internal/handler"
	"drk-url-shortener/internal/lib/logger/sl"
	"drk-url-shortener/internal/repository"
	"drk-url-shortener/internal/service"
	"log/slog"
	"net/http"
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

	// 📌Stabilization Stage (Production-Ready MVP)
	// 3️⃣ Инфраструктурный слой- db
	// Будет осуществлено в iter9 (сохранение сокращенных URL в файл при выходе и загрузка из него при запуске)) и затем в iter10 (pg)

	// DI
	repos := repository.New(make(map[string]string))
	services := service.New(repos)
	handlers := handler.New(services, cfg, log)

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
