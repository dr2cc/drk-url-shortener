package v1

import (
	"drk-url-shortener/internal/usecase/mocks"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRouter_shortenText(t *testing.T) {
	// Мы тестируем Router, поэтому он — sut (Система Под Тестом)
	// 1. Передаем мок напрямую в роутер.
	// 2. Вызываем метод у sut (здесь shortenText).
	// 3. Проверяем, как хендлер отреагирует на ответы от интерактора Shortener.
	// (interactor- «тот, кто управляет взаимодействием»).

	// mockBehavior (имитация поведения), тип-функция (function type), настройщик поведения мока.
	// Это callback-функция (так как эта логика передается внутрь теста, чтобы сработать в нужный момент), инъекция поведения.
	// В данном случае принимает объект (структуру) имитирующий UseCase interface и ...

	// // В поведение передаем мок UseCase и три основных параметра
	// type mockBehavior func(ucMock *mocks.MockUseCase, url string, baseURL string, slug string)

	type testCase struct {
		name               string
		url                string
		baseURL            string
		slug               string
		expectedStatusCode int
		expectedBody       string
		body               io.ReadCloser
		mockBehavior       func(ucMock *mocks.MockUseCase, tc *testCase) // Передаем указатель на себя
	}

	// 2. Описываем таблицу тестов
	tests := []testCase{
		{
			name:               "OK",
			url:                "https://google.com",
			baseURL:            "http://localhost:8080",
			slug:               "abc",
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "http://localhost:8080/abc",
			body:               io.NopCloser(strings.NewReader("https://google.com")),
			mockBehavior: func(ucMock *mocks.MockUseCase, tc *testCase) {
				// «Когда роутер вызовет метод Shorten("https://google.com"), ничего не пиши в БД,
				// а сразу Верни строку "abc" и ошибку nil»
				ucMock.EXPECT().Shorten(tc.url).Return(tc.slug, nil)
				// «Когда роутер вызовет метод FormatShortURL("http://localhost:8080","abc"), ничего не делай,
				// а сразу Верни строку "http://localhost:8080/abc"»
				ucMock.EXPECT().FormatShortURL(tc.baseURL, tc.slug).Return(tc.expectedBody)
			},
		},
		{
			// Ошибка чтения из body
			name:               "Scroll reading error",
			url:                "",
			baseURL:            "http://localhost:8080",
			slug:               "",
			expectedStatusCode: http.StatusBadRequest,
			// Создаем прямо в строке таблицы, ничего заранее описывать не нужно:
			body:         io.NopCloser(iotest.ErrReader(errors.New("read error"))),
			mockBehavior: func(ucMock *mocks.MockUseCase, tc *testCase) {}, // Хендлер упадет до UseCase
		},
		{
			name:               "Empty URL",
			url:                "",
			baseURL:            "http://localhost:8080",
			slug:               "",
			expectedStatusCode: http.StatusBadRequest,
			body:               io.NopCloser(strings.NewReader("")),
			mockBehavior:       func(ucMock *mocks.MockUseCase, tc *testCase) {}, // Хендлер отбракует запрос до UseCase
		},
		{
			name:               "SaveURL Error",
			url:                "https://google.com",
			baseURL:            "http://localhost:8080",
			slug:               "",
			expectedStatusCode: http.StatusBadRequest,
			body:               io.NopCloser(strings.NewReader("https://google.com")),
			mockBehavior: func(ucMock *mocks.MockUseCase, tc *testCase) {
				// ... верни любую ошибку
				ucMock.EXPECT().Shorten(tc.url).Return(tc.slug, errors.New("internal server error"))
				// до второго метода не дойдет!
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
			tt.mockBehavior(ucMock, &tt)

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
			req := httptest.NewRequest("POST", "/", tt.body) //bytes.NewReader([]byte(tt.url)))

			// Выполнение действия (Act)
			r.ServeHTTP(w, req)

			// Проверка утверждений (Assert)
			// 3. Проверяем, как хендлер отреагирует на ответы от интерактора Shortener.
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			if tt.name == "OK" {
				assert.Equal(t, "http://localhost:8080/abc", w.Body.String())
			}
		})
	}
}
