package v1

import (
	"drk-url-shortener/internal/usecase"
	"drk-url-shortener/internal/usecase/mocks"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_router_redirect(t *testing.T) {
	// Тестируется ЗАВИСИМОСТИ!
	// Проверяем, как хендлер отреагирует на ответы от интерактора (ex. сервиса) Shortener.
	// Interactor означает «тот, кто управляет взаимодействием».

	// mockBehavior (имитация поведения), тип-функция (function type), настройщик поведения мока.
	// Это callback-функция (так как эта логика передается внутрь теста, чтобы сработать в нужный момент), инъекция поведения.
	// В данном случае принимает объект (структуру) имитирующий UseCase interface и строку слага.

	// В поведение передаем мок UseCase
	type mockBehavior func(ucMock *mocks.MockUseCase, slug string)

	tests := []struct {
		name               string
		id                 string
		expectedStatusCode int
		expectedBody       string
		mockBehavior       mockBehavior
	}{
		{
			name:               "OK",
			id:                 "abc",
			expectedStatusCode: http.StatusTemporaryRedirect,
			expectedBody:       "", // при редиректе тело обычно не проверяем
			mockBehavior: func(ucMock *mocks.MockUseCase, slug string) {
				// «Когда роутер вызовет метод GetOriginal("abc"), ничего не ищи в БД, а сразу верни строку "https://google.com" и ошибку nil . Повтори один раз.».
				// Метод EXPECT() есть у каждого сгенерированного мока.
				// Он возвращает специальный объект-регистратор (recorder *MockUseCaseMockRecorder — указатель на "записывающий" объект).
				// Задача этого объекта — записывать, какие методы должен вызвать ваш код во время теста.
				ucMock.EXPECT(). // "При обращении к объекту ucMock мы будем ОЖИДАТЬ()"
							GetOriginal(slug). // метод GetOriginal вызывается не у самого мока, а у регистратора, которого вернул нам s.EXPECT().
					// Внутри сгенерированного кода этот метод создает структуру Call.
					// Эта структура запоминает, какие аргументы (slug) ожидается получить.
					// Метод возвращает эту самую структуру Call.
					Return("https://google.com", nil). // Поскольку Get вернул объект Call, мы можем вызывать его методы через точку.
					// Return записывает в объект Call, что именно нужно вернуть при вызове.
					// Times записывает, сколько раз этот метод должен быть вызван.
					// Каждый из этих методов снова возвращает тот же самый объект Call. Это и позволяет писать их цепочкой друг за другом.
					// В данном случае,
					// как только программа вызовет метод GetOriginal,
					// имитатор мгновенно отдаст ей "https://google.com" (url) и nil (отсутствие ошибки).
					// Это позволяет тестировать логику дальше, не обращаясь к реальной базе данных.
					Times(1) // (не обязательно) - сколько раз вызываем (по умолчанию 1)
				//
				// Когда функция mockBehavior выполнится, GoMock создаст объект Call со следующими значениями:
				// - method: "GetOriginal" — GoMock запомнил, какой метод мы ждем.
				// - argumentCheck: Сюда запишется матчер (проверяльщик),
				// который жестко сравнивает входящую строку с переменной slug. Он сработает как gomock.Eq(slug).
				// - rets: Слайс из двух элементов: ["https://google.com", nil]. Их мок отдаст обратно вашему коду.
				// ...
			},
		},
		{
			name:               "empty db",
			id:                 "abc",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "short link not found",
			// «Когда роутер вызовет метод GetOriginal("abc"), притворись, что в БД ничего нет, и верни пустую строку и ошибку ErrCodeNotFound . Повтори один раз.».
			mockBehavior: func(ucMock *mocks.MockUseCase, slug string) {
				ucMock.EXPECT().
					GetOriginal(slug).
					Return("", usecase.ErrCodeNotFound).
					Times(1)
			},
		},
		{
			name:               "internal db error",
			id:                 "abc",
			expectedStatusCode: http.StatusBadRequest, // по ТЗ
			expectedBody:       "internal server error",
			mockBehavior: func(ucMock *mocks.MockUseCase, slug string) {
				ucMock.EXPECT().
					GetOriginal(slug).
					Return("", errors.New("internal server error")). // Любая другая ошибка (например, упала БД)
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Инициализация моков.
			ctrl := gomock.NewController(t)

			// Создаем ложный юзкейс
			ucMock := mocks.NewMockUseCase(ctrl)

			// Настраиваем поведение мока под конкретный тест-кейс
			tt.mockBehavior(ucMock, tt.id)
			// tt.mockBehavior должен принять параметры- ucMock *mocks.MockUseCase и slug
			// Он принимает ucMock созданный выше и задает ему условия работы данного тест-кейса:
			// ucMock.EXPECT().GetOriginal(tt.id).Return("https://google.com", nil).Times(1)

			// 2. Создаем логгер-заглушку, который пишет "в никуда".
			// Без него было можно работать, пока не появился негативный тест-кейс.
			discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

			// 3. Создаем Систему Под Тестом (SUT)
			// Мы тестируем именно Router, поэтому он — sut!
			// Передаем мок напрямую в роутер
			sut := &Router{
				shortener: ucMock, // Роутер зависит от интерфейса usecase.UseCase
				// MockUseCase struct реализует интерфейс usecase.UseCase !
				log: discardLogger,
			}

			// 4. Настройка окружения (Инфраструктура HTTP)
			r := chi.NewRouter()
			r.Get("/{id}", sut.redirect) // вызываем метод у sut

			w := httptest.NewRecorder()
			// Запрос всегда правильный ("/abc").
			req := httptest.NewRequest("GET", "/"+tt.id, nil)

			// 5. Выполнение действия (Act)
			r.ServeHTTP(w, req)

			// 6. Проверка утверждений (Assert)
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Для успешного теста (OK)
			if tt.expectedStatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
				// w.Header().Get("Location") достает адрес, куда хендлер делает редирект
				assert.Equal(t, "https://google.com", w.Header().Get("Location"))
			} else {
				// Для негативных тестов (например, 400 Bad Request)
				// Проверяем, что клиенту возвращается вменяемый текст ошибки
				// (подставьте сюда вашу логику: JSON или обычная строка, например "code not found")
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}
