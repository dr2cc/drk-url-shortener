package repository

import (
	"drk-url-shortener/internal/usecase"
	"drk-url-shortener/internal/usecase/mocks"
	"errors"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCache_Save(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		wantSlug     string
		wantErr      error
		mockBehavior func(repoMock *mocks.MockRepository, genMock *mocks.MockCodeGenerator)
	}{
		{
			name:     "OK",
			url:      "https://google.com",
			wantSlug: "abc",
			wantErr:  nil,
			mockBehavior: func(repoMock *mocks.MockRepository, genMock *mocks.MockCodeGenerator) {
				// 1. Сначала логика просит генератор создать код "abc"
				genMock.EXPECT().RandomString().Return("abc")
				// 2. «Когда (кто?) вызовет метод Save(,"abc","https://google.com"), ничего не пиши в БД,
				// а сразу Верни ошибку nil»
				repoMock.EXPECT().Save("abc", "https://google.com").Return(nil)
			},
		},
		{
			name:     "Repository Error",
			url:      "https://google.com",
			wantSlug: "",
			wantErr:  errors.New("db connection failure"),
			mockBehavior: func(repoMock *mocks.MockRepository, genMock *mocks.MockCodeGenerator) {
				genMock.EXPECT().RandomString().Return("abc")
				// Симулируем падение базы данных
				repoMock.EXPECT().Save("abc", "https://google.com").Return(errors.New("db connection failure"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// // ТРИ пункта обязательные в каждом тесте созданном по этой схеме.
			// 1. Инициализация моков.
			ctrl := gomock.NewController(t)

			// 2. Создание имитатора.
			// repoMock должен реализовать методы usecase.Repository
			repoMock := mocks.NewMockRepository(ctrl)
			genMock := mocks.NewMockCodeGenerator(ctrl)

			// 3. Настройка поведения мока под конкретный тест-кейс
			tt.mockBehavior(repoMock, genMock)

			// 4. Инициализируем тестируемый слой бизнес-логики (замените на ваш конструктор)
			service := usecase.New(repoMock, genMock)

			// 5. Вызываем тестируемый метод
			// Вызываем целевой метод
			slug, err := service.Shorten(tt.url)

			// Проверка ошибок через require
			if tt.wantErr != nil {
				require.Error(t, err)
				// Сравниваем строки ошибок, если это динамическая ошибка
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}

			// Проверка возвращенного slug через assert
			assert.Equal(t, tt.wantSlug, slug)
		})
	}
}
