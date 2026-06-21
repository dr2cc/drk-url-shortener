package v1

import (
	"drk-url-shortener/internal/lib/testlog"
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
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const baseURL = "http://localhost:8080"

func TestRouter_shortenJSON(t *testing.T) {

	tests := []struct {
		name               string
		expectedStatusCode int
		expectedBody       string
		url                string
		isRawInput         bool   // Флаг: использовать ли "сырой" input вместо генерации через fmt.Sprintf
		rawInput           string // Сам "сырой" текст для передачи битого JSON или пустоты
		mockBehavior       func(ucMock *mocks.MockUseCase)
	}{
		{
			name:               "OK",
			expectedStatusCode: http.StatusCreated,
			expectedBody:       `{"result": "http://localhost:8080/abc"}`,
			url:                "https://google.com",
			mockBehavior: func(ucMock *mocks.MockUseCase) {
				// «Когда роутер вызовет метод Shorten("https://google.com"), ничего не пиши в БД,
				// а сразу Верни строку "abc" и ошибку nil»
				ucMock.EXPECT().Shorten("https://google.com").Return("abc", nil)
				// «Когда роутер вызовет метод FormatShortURL("http://localhost:8080","abc"), ничего не делай,
				// а сразу Верни строку "http://localhost:8080/abc"»
				ucMock.EXPECT().FormatShortURL("http://localhost:8080", "abc").Return("http://localhost:8080/abc")
			},
		},
		{
			name:               "Invalid JSON",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "failed to decode request body", // То, что возвращает ваш хендлер при ошибке десериализации
			isRawInput:         true,
			rawInput:           `{"url": "https://google.com"`, // Сломанный JSON (нет закрывающей скобки)
			mockBehavior:       func(ucMock *mocks.MockUseCase) {},
		},
		{
			name:               "Empty Body",
			expectedStatusCode: http.StatusBadRequest,
			mockBehavior:       func(ucMock *mocks.MockUseCase) {},
		},
		{
			name:               "Validation Error - Empty URL",
			expectedStatusCode: http.StatusBadRequest,
			url:                "",
			mockBehavior:       func(ucMock *mocks.MockUseCase) {},
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
			tt.mockBehavior(ucMock)

			// Логгер будет писать в стандартный механизм тестов Go.
			// Вы увидите логи в консоли ТОЛЬКО если тест завершился ошибкой (go test -v).
			// Главный плюс- будет виден весь лог до момента ошибки.
			log := testlog.New(t)

			// Система Под Тестом (SUT)
			// Передаем мок напрямую в роутер
			sut := &Router{
				shortener: ucMock,
				baseURL:   baseURL,
				validator: validator.New(),
				log:       log,
			}

			// Настройка окружения (Инфраструктура HTTP)
			r := chi.NewRouter()
			// Вызываем метод у sut
			r.Post("/api/shorten", sut.shortenJSON)

			w := httptest.NewRecorder()

			// Адаптивное формирование тела запроса
			var input string
			if tt.isRawInput {
				input = tt.rawInput
			} else {
				input = fmt.Sprintf(`{"url": "%s"}`, tt.url)
			}

			// Передаем strings.NewReader напрямую. httptest сам обернет его в io.ReadCloser
			// и корректно посчитает длину тела (ContentLength).
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(input))

			// Добавляем заголовок, чтобы обработчик понял, что это JSON
			req.Header.Set("Content-Type", "application/json")

			// // Выполнение действия (Act)
			// r.ServeHTTP(w, req)
			// Выполнение действия (Act) с красивым перехватом паники
			assert.NotPanics(t, func() {
				r.ServeHTTP(w, req)
			}, "The handler panicked! Check the initialization of dependencies.")

			// Проверка утверждений (Assert)
			// Проверяем, как хендлер отреагирует на ответы от интерактора Shortener.
			// Тест прервется сразу же на этой строчке, если статус не совпадет
			require.Equal(t, tt.expectedStatusCode, w.Code, "Invalid status code. Response: %s", w.Body.String())

			if tt.name == "OK" {
				// Для успешного кейса идеально подходит JSONEq (он проигнорирует пробелы и \n)
				assert.JSONEq(t, tt.expectedBody, w.Body.String(), "The handler's response does not match the expected JSON template.")
			} else {
				// Для ошибок (Plain Text) используем assert.Contains.
				// Он проверяет, что строка tt.expectedBody есть внутри ответа,
				// и ему абсолютно плевать на автоматический перевод строки \n в конце!
				assert.Contains(t, w.Body.String(), tt.expectedBody, "The error message in the response is incorrect.")
			}
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

	// // В поведение передаем мок UseCase
	// type mockBehavior func(ucMock *mocks.MockUseCase)

	tests := []struct {
		name               string
		expectedStatusCode int
		expectedBody       string
		body               io.ReadCloser
		mockBehavior       func(ucMock *mocks.MockUseCase)
	}{
		{
			name:               "OK",
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "http://localhost:8080/abc", // w.Write([]byte(content)) не добавляет \n
			body:               io.NopCloser(strings.NewReader("https://google.com")),
			mockBehavior: func(ucMock *mocks.MockUseCase) {
				// «Когда роутер вызовет метод Shorten("https://google.com"), ничего не пиши в БД,
				// а сразу Верни строку "abc" и ошибку nil»
				ucMock.EXPECT().Shorten("https://google.com").Return("abc", nil)
				// «Когда роутер вызовет метод FormatShortURL("http://localhost:8080","abc"), ничего не делай,
				// а сразу Верни строку "http://localhost:8080/abc"»
				ucMock.EXPECT().FormatShortURL("http://localhost:8080", "abc").Return("http://localhost:8080/abc")
			},
		},
		{
			name:               "Error reading from body",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "scroll reading error\n", // http.Error добавляет \n
			// Создаем прямо в строке таблицы, ничего заранее описывать не нужно:
			body:         io.NopCloser(iotest.ErrReader(errors.New("read error"))),
			mockBehavior: func(ucMock *mocks.MockUseCase) {}, // Хендлер упадет до UseCase
		},
		{
			name:               "Empty URL",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "URL not found in request body\n", // http.Error добавляет \n
			body:               io.NopCloser(strings.NewReader("")),
			mockBehavior:       func(ucMock *mocks.MockUseCase) {}, // Хендлер отбракует запрос до UseCase
		},
		{
			name:               "SaveURL Error",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "failed to add url\n", // http.Error добавляет \n
			body:               io.NopCloser(strings.NewReader("https://google.com")),
			mockBehavior: func(ucMock *mocks.MockUseCase) {
				// ... верни пустую строку и любую ошибку
				ucMock.EXPECT().Shorten("https://google.com").Return("", errors.New("internal server error"))
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
			tt.mockBehavior(ucMock)

			// Логгер-заглушка.
			discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

			// Система Под Тестом (SUT)
			// 1. Передаем мок напрямую в роутер
			sut := &Router{
				shortener: ucMock,
				baseURL:   baseURL,
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
