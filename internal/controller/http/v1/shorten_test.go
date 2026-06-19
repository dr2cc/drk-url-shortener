package v1

import (
	"drk-url-shortener/internal/usecase/mocks"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRouter_shortenJSON(t *testing.T) {
	type testCase struct {
		name               string
		baseURL            string
		slug               string
		expectedStatusCode int
		expectedBody       string
		url                string
		// body         *bytes.Reader
		mockBehavior func(ucMock *mocks.MockUseCase, tc testCase)
	}

	tests := []testCase{
		{
			name:               "OK",
			baseURL:            "http://localhost:8080",
			slug:               "abc",
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "http://localhost:8080/abc", // w.Write([]byte(content)) не добавляет \n
			url:                "https://google.com",
			// body: bytes.NewReader([]byte(fmt.Sprintf(`{"url": "%s"}`, "https://google.com"))),
			// body: bytes.NewReader(fmt.Appendf(nil, `{"url": "%s"}`, "https://google.com")),
			mockBehavior: func(ucMock *mocks.MockUseCase, tc testCase) {
				// «Когда роутер вызовет метод Shorten("https://google.com"), ничего не пиши в БД,
				// а сразу Верни строку "abc" и ошибку nil»
				ucMock.EXPECT().Shorten("https://google.com").Return(tc.slug, nil)
				// «Когда роутер вызовет метод FormatShortURL("http://localhost:8080","abc"), ничего не делай,
				// а сразу Верни строку "http://localhost:8080/abc"»
				ucMock.EXPECT().FormatShortURL(tc.baseURL, tc.slug).Return(tc.expectedBody)
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
			tt.mockBehavior(ucMock, tt)

			// Логгер-заглушка.
			discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

			// Система Под Тестом (SUT)
			// 1. Передаем мок напрямую в роутер
			sut := &Router{
				shortener: ucMock,
				baseURL:   tt.baseURL,
				validator: validator.New(),
				log:       discardLogger,
			}

			// Настройка окружения (Инфраструктура HTTP)
			r := chi.NewRouter()
			// 2. Вызываем метод у sut
			r.Post("/api/shorten", sut.shortenJSON)

			w := httptest.NewRecorder()
			// input := fmt.Sprintf(`{"url": "%s"}`, tt.url)
			// req := httptest.NewRequest("POST", "/api/shorten", tt.body) //bytes.NewReader([]byte(tt.url)))

			input := fmt.Sprintf(`{"url": "%s"}`, tt.url)

			// Передаем strings.NewReader напрямую. httptest сам обернет его в io.ReadCloser
			// и корректно посчитает длину тела (ContentLength).
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(input))

			// Не забываем добавить заголовок, чтобы обработчик понял, что это JSON
			req.Header.Set("Content-Type", "application/json")

			// Выполнение действия (Act)
			r.ServeHTTP(w, req)

			// Проверка утверждений (Assert)
			// 3. Проверяем, как хендлер отреагирует на ответы от интерактора Shortener.
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			// assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}

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
		baseURL            string
		slug               string
		expectedStatusCode int
		expectedBody       string
		body               io.ReadCloser
		mockBehavior       func(ucMock *mocks.MockUseCase, tc testCase)
	}

	tests := []testCase{
		{
			name:               "OK",
			baseURL:            "http://localhost:8080",
			slug:               "abc",
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "http://localhost:8080/abc", // w.Write([]byte(content)) не добавляет \n
			body:               io.NopCloser(strings.NewReader("https://google.com")),
			mockBehavior: func(ucMock *mocks.MockUseCase, tc testCase) {
				// «Когда роутер вызовет метод Shorten("https://google.com"), ничего не пиши в БД,
				// а сразу Верни строку "abc" и ошибку nil»
				ucMock.EXPECT().Shorten("https://google.com").Return(tc.slug, nil)
				// «Когда роутер вызовет метод FormatShortURL("http://localhost:8080","abc"), ничего не делай,
				// а сразу Верни строку "http://localhost:8080/abc"»
				ucMock.EXPECT().FormatShortURL(tc.baseURL, tc.slug).Return(tc.expectedBody)
			},
		},
		{
			name:               "Error reading from body",
			baseURL:            "http://localhost:8080",
			slug:               "",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "scroll reading error\n", // http.Error добавляет \n
			// Создаем прямо в строке таблицы, ничего заранее описывать не нужно:
			body:         io.NopCloser(iotest.ErrReader(errors.New("read error"))),
			mockBehavior: func(ucMock *mocks.MockUseCase, tc testCase) {}, // Хендлер упадет до UseCase
		},
		{
			name:               "Empty URL",
			baseURL:            "http://localhost:8080",
			slug:               "",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "URL not found in request body\n", // http.Error добавляет \n
			body:               io.NopCloser(strings.NewReader("")),
			mockBehavior:       func(ucMock *mocks.MockUseCase, tc testCase) {}, // Хендлер отбракует запрос до UseCase
		},
		{
			name:               "SaveURL Error",
			baseURL:            "http://localhost:8080",
			slug:               "",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "failed to add url\n", // http.Error добавляет \n
			body:               io.NopCloser(strings.NewReader("https://google.com")),
			mockBehavior: func(ucMock *mocks.MockUseCase, tc testCase) {
				// ... верни любую ошибку
				ucMock.EXPECT().Shorten("https://google.com").Return(tc.slug, errors.New("internal server error"))
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
			tt.mockBehavior(ucMock, tt)

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
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
