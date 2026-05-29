package v1

import (
	"drk-url-shortener/internal/repository"
	"drk-url-shortener/internal/usecase"
	"drk-url-shortener/internal/usecase/mocks"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func Test_router_redirect(t *testing.T) {
	// Тестируется ЗАВИСИМОСТИ!
	// Проверяем, как хендлер отреагирует на ответы от базы данных.

	// mockBehavior (имитация поведения), тип-функция (function type), настройщик поведения мока.
	// Callback-функция — так как эта логика передается внутрь теста, чтобы сработать в нужный момент (инъекция поведения).
	// В данном случае принимает объект (структуру) имитирующий ShortURLRepo interface и строку слага.
	type mockBehavior func(s *mocks.MockShortURLRepo, slug string)

	tests := []struct {
		name               string
		id                 string
		expectedStatusCode int
		mockBehavior       mockBehavior
	}{
		{
			name:               "OK",
			id:                 "abc",
			expectedStatusCode: http.StatusTemporaryRedirect,
			mockBehavior: func(s *mocks.MockShortURLRepo, slug string) {
				// Метод EXPECT() есть у каждого сгенерированного мока.
				// Он возвращает специальный объект-регистратор (recorder *MockShortURLMockRecorder — указатель на "записывающий" объект).
				// Задача этого объекта — записывать, какие методы должен вызвать ваш код во время теста.
				// Сообщаем моку: «Сейчас я опишу вызов, который должен произойти во время работы программы» или
				s.EXPECT(). // "При обращении к объекту s мы будем ОЖИДАТЬ()"
						Get(slug). // метод Get вызывается не у самого мока, а у регистратора, которого вернул нам s.EXPECT().
					// Внутри сгенерированного кода этот метод создает структуру Call.
					// Эта структура запоминает, какие аргументы (slug) ожидается получить.
					// Метод возвращает эту самую структуру Call.
					Return("https://google.com", nil). // Поскольку Get вернул объект Call, мы можем вызывать его методы через точку.
					// Return записывает в объект Call, что именно нужно вернуть при вызове.
					// Times записывает, сколько раз этот метод должен быть вызван.
					// Каждый из этих методов снова возвращает тот же самый объект Call. Это и позволяет писать их цепочкой друг за другом.
					// В данном случае,
					// как только программа вызовет метод Get,
					// имитатор мгновенно отдаст ей "https://google.com" (url) и nil (отсутствие ошибки).
					// Это позволяет тестировать логику дальше, не обращаясь к реальной базе данных.
					Times(1) // (не обязательно) - сколько раз вызываем (по умолчанию 1)
				//
				// Когда функция mockBehavior выполнится, GoMock создаст объект Call со следующими значениями:
				// - method: "Get" — GoMock запомнил, какой метод мы ждем.
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
			mockBehavior: func(s *mocks.MockShortURLRepo, slug string) {
				s.EXPECT().
					Get(slug).
					Return("", repository.ErrNotFound).
					Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализация моков.
			ctrl := gomock.NewController(t)
			// Создаем "ложный" сервис, который "притворяется" реальной бизнес-логикой (интерфейсом ShortURLRepo).
			repo := mocks.NewMockShortURLRepo(ctrl)
			// Передавая параметры (repo, tt.id) мы указываем, что
			// 🕖 ожидаем получить вызов методов сервиса repo, а в качестве аргумента передадим id
			tt.mockBehavior(repo, tt.id)
			// Создаем объект сервисов, но передадим аргументом для интерфейса ShortURLRepo наш "ложный" repo.
			// usecase.Shortener "думает", что работает с настоящей базой или API, хотя на самом деле он работает с контролируемым нами моком.
			services := &usecase.Shortener{Repo: repo}

			// 1. Инициализируем хендлер.
			// Структура Router получает объект services, внутри которого уже есть наш мок.
			// Хендлер не знает, как получить url, он лишь делегирует это сервису.
			// Создаем логгер-заглушку, который пишет "в никуда".
			// Без него было можно работать, пока не появился негативный тест-кейс.
			discardLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := &Router{
				shortener: services,
				log:       discardLogger,
			}

			// 2. Init Endpoint
			r := chi.NewRouter()
			// Регистрируем конкретную функцию контроллера (тестируемый метод хендлера redirect) на маршрут /{id}
			// chi теперь знает: «Если придет GET-запрос на этот адрес, нужно запустить именно этот код».
			// А при вызове handler.redirect внутри сработает цепочка, ведущая к моку.
			r.Get("/{id}", handler.redirect)

			w := httptest.NewRecorder()
			// Запрос всегда правильный ("/abc"). Все негативные сценарии- в db (мок!)
			req := httptest.NewRequest("GET", "/abc", nil)

			// Make Request
			// Вызываем метод ServeHTTP у объекта роутера
			// Отдаем роутеру «виртуальный» запрос (req) и «записывающее устройство» (w).
			// Роутер прогоняет запрос через свои механизмы, вызывает хендлер, тот вызывает сервис (мок!),
			// получает ответ и записывает результат в w.
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Для успешного теста (OK)
			if tt.expectedStatusCode == http.StatusTemporaryRedirect {
				assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
				// w.Header().Get("Location") достает адрес, куда хендлер делает редирект
				assert.Equal(t, "https://google.com", w.Header().Get("Location"))
			}
		})
	}
}
