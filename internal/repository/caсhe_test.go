package repository

import (
	"drk-url-shortener/internal/usecase/mocks"
	"testing"

	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestCache_Save(t *testing.T) {
	tests := []struct {
		name string
		// alias        string
		// url          string
		err          error
		mockBehavior func(repoMock *mocks.MockRepository)
	}{
		{
			name: "OK",
			// alias: "abc",
			// url:   "https://google.com",
			err: nil,
			mockBehavior: func(repoMock *mocks.MockRepository) {
				// «Когда (кто?) вызовет метод Save("https://google.com","abc"), ничего не пиши в БД,
				// а сразу Верни ошибку nil»
				repoMock.EXPECT().Save("https://google.com", "abc").Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// // ТРИ пункта обязательные в каждом тесте созданном по этой схеме.
			// 1. Инициализация моков.
			ctrl := gomock.NewController(t)

			// 2. Создаем имитатор.
			// repoMock должен реализовать методы usecase.Repository
			repoMock := mocks.NewMockRepository(ctrl)

			// 3. Настраиваем поведение мока под конкретный тест-кейс
			tt.mockBehavior(repoMock)

			// К 01.07.26 - Что тестируем здесь (Система Под Тестом (SUT))?
			// Видимо основную структуру?
			// Сдался, спросил: https://share.google/aimode/J4WjGSMbED1ky7nIF

			// // Логгер будет писать в стандартный механизм тестов Go.
			// // Вы увидите логи в консоли ТОЛЬКО если тест завершился ошибкой (go test -v).
			// // Главный плюс- будет виден весь лог до момента ошибки.
			// log := testlog.New(t)

			// // Система Под Тестом (SUT)
			// // Передаем мок напрямую в роутер
			// sut := &Router{
			// 	shortener: ucMock,
			// 	baseURL:   baseURL,
			// 	log:       log,
			// }

			// // Настройка окружения (Инфраструктура HTTP)
			// r := chi.NewRouter()
			// // Вызываем метод у sut
			// r.Post("/api/shorten", sut.shortenJSON)

			// w := httptest.NewRecorder()

			// // Передаем strings.NewReader напрямую. httptest сам обернет его в io.ReadCloser
			// // и корректно посчитает длину тела (ContentLength).
			// req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.rawInputBody))

			// // Добавляем заголовок, чтобы обработчик понял, что это JSON
			// req.Header.Set("Content-Type", "application/json")

			// Проверка утверждений (Assert)
			// Третим аргументом должен быть результат теста (сейчас заглушка)!
			assert.Equal(t, tt.err, nil)

		})
	}
}
