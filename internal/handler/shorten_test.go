package handler

import (
	"bytes"
	"drk-url-shortener/internal/config"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestShortenTextHandler(t *testing.T) {
	// 1️⃣ изменение- СОЗДАЕМ новый конфиг (копируем код в main)
	// 1. Подготовка общих зависимостей (фикстур)
	cfg := config.Config{BaseURL: "http://localhost:8080"}

	// Описываем структуру тест-кейса
	type testCase struct {
		name           string
		method         string
		body           string
		wantStatusCode int
		// wantBodyCheck — функция для дополнительной проверки содержимого ответа (body) и состояния БД (repo)
		wantBodyCheck func(t *testing.T, body string, repo map[string]string)
	}

	// 2. Определение сценариев тестирования
	tests := []testCase{
		{
			name:           "Успешное создание короткой ссылки",
			method:         http.MethodPost,
			body:           "example.com",
			wantStatusCode: http.StatusCreated,
			wantBodyCheck: func(t *testing.T, body string, repo map[string]string) {
				// Проверяем, что ответ начинается с BaseURL
				if !strings.HasPrefix(body, cfg.BaseURL+"/") {
					t.Errorf("неверный формат ответа: %s", body)
				}
				// Проверяем, что в базе (map) появилась запись
				if len(repo) != 1 {
					t.Errorf("ожидалась 1 запись в репозитории, найдено: %d", len(repo))
				}
			},
		},
		{
			name:   "Ошибка: неверный HTTP метод (GET)",
			method: http.MethodGet,
			body:   "https://example.com",
			// В коде теста сделана привязка r.Post("/", shortenText(...)).
			// Когда отправляется запрос методом GET на адрес /,
			// роутер chi понимает, что данный маршрут существует, но не поддерживает метод GET.
			// Он перехватывает запрос и возвращает статус 405, даже не доходя до выполнения функции shortenText
			wantStatusCode: http.StatusMethodNotAllowed, // До 400 недойдет, ожидаем 405
			wantBodyCheck: func(t *testing.T, body string, repo map[string]string) {
				// Роутер chi по умолчанию возвращает пустой body или стандартный текст для 405 ошибки.
				// Проверяем, что в репозиторий ничего не записалось
				if len(repo) > 0 {
					t.Errorf("репозиторий должен быть пустым")
				}
			},
		},
		{
			name:           "Ошибка: пустой body",
			method:         http.MethodPost,
			body:           "",
			wantStatusCode: http.StatusBadRequest,
			wantBodyCheck: func(t *testing.T, body string, repo map[string]string) {
				if !strings.Contains(body, "URL not found in request body") {
					t.Errorf("неожиданный текст ошибки: %s", body)
				}
			},
		},
	}

	// 3. Итерация и запуск каждого сценария
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 2️⃣ изменение- СОЗДАЕМ хранилище!
			// Изолированный репозиторий для каждого подтеста
			repo := make(map[string]string)

			// Создаем роутер через вашу функцию New
			// Для теста логгер можно передать как nil или заглушку, если chi-slog это позволяет,
			// либо инициализировать slog.New(slog.NewTextHandler(io.Discard, nil))
			r := chi.NewRouter()
			r.Post("/", shortenText(repo, cfg))

			// Создаем виртуальный запрос и рекордер ответа
			//
			// Генерирует объект *http.Request без поднятия сетевых сокетов. Данные тела передаются через bytes.NewBufferString.
			req := httptest.NewRequest(tc.method, "/", bytes.NewBufferString(tc.body))
			// Действует как браузер. Записывает заголовки, статус-код и тело ответа, которые возвращает обработчик.
			rr := httptest.NewRecorder()

			// Выполняем запрос через роутер
			r.ServeHTTP(rr, req)

			// Проверяем статус-код
			if rr.Code != tc.wantStatusCode {
				t.Errorf("код ответа: получили %d, ожидали %d", rr.Code, tc.wantStatusCode)
			}

			// Читаем тело ответа
			respBody, _ := io.ReadAll(rr.Body)
			gotBody := string(bytes.TrimSpace(respBody))

			// Вызываем специфичные для кейса проверки
			if tc.wantBodyCheck != nil {
				tc.wantBodyCheck(t, gotBody, repo)
			}
		})
	}
}
