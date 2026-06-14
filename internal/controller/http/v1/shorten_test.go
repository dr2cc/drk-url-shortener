package v1

import (
	"bytes"
	"drk-url-shortener/internal/usecase/mocks"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRouter_shortenText(t *testing.T) {
	// Тестируется ЗАВИСИМОСТИ!

	// Мы тестируем именно Router, поэтому он — sut (Система Под Тестом)
	// 1. Передаем мок напрямую в роутер.
	// 2. Вызываем метод у sut (здесь shortenText).
	// 3. Проверяем, как хендлер отреагирует на ответы от интерактора Shortener.
	// (interactor- «тот, кто управляет взаимодействием»).

	// mockBehavior (имитация поведения), тип-функция (function type), настройщик поведения мока.
	// Это callback-функция (так как эта логика передается внутрь теста, чтобы сработать в нужный момент), инъекция поведения.
	// В данном случае принимает объект (структуру) имитирующий UseCase interface и ...

	// В поведение передаем мок UseCase
	type mockBehavior func(ucMock *mocks.MockUseCase, url string, baseURL string, slug string, expectedURL string)

	tests := []struct {
		name               string
		baseURL            string
		url                string
		slug               string
		expectedURL        string
		expectedStatusCode int
		mockBehavior       mockBehavior
	}{
		{
			name:               "OK",
			baseURL:            "http://localhost:8080",
			url:                "https://google.com",
			slug:               "abc",
			expectedURL:        "http://localhost:8080/abc",
			expectedStatusCode: http.StatusCreated,
			mockBehavior: func(ucMock *mocks.MockUseCase, url string, baseURL string, slug string, expectedURL string) {
				// «Когда роутер вызовет метод Shorten("https://google.com"), ничего не пиши в БД,
				// а сразу Верни строку "abc" и ошибку nil»
				ucMock.EXPECT().Shorten(url).Return(slug, nil)
				// «Когда роутер вызовет метод FormatShortURL("http://localhost:8080","abc"), ничего не делай,
				// а сразу Верни строку "http://localhost:8080/abc"»
				ucMock.EXPECT().FormatShortURL(baseURL, slug).Return(expectedURL)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализация моков.
			ctrl := gomock.NewController(t)

			// Создаем ложный юзкейс
			// ucMock реализует методы usecase.UseCase
			ucMock := mocks.NewMockUseCase(ctrl)

			// Настраиваем поведение мока под конкретный тест-кейс
			tt.mockBehavior(ucMock, tt.url, tt.baseURL, tt.slug, tt.expectedURL)

			// Логгер-заглушка.
			discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

			// Система Под Тестом (SUT)
			// 1. Передаем мок напрямую в роутер
			sut := &Router{
				shortener: ucMock,
				baseURL:   tt.baseURL,
				log:       discardLogger,
			}

			// Настройка окружения (Инфраструктура HTTP)
			r := chi.NewRouter()
			// 2. Вызываем метод у sut
			r.Post("/", sut.shortenText)

			w := httptest.NewRecorder()
			//
			req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(tt.url)))

			// Выполнение действия (Act)
			r.ServeHTTP(w, req)

			// Проверка утверждений (Assert)
			// 3. Проверяем, как хендлер отреагирует на ответы от интерактора Shortener.
			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
