package v1

import (
	"testing"
)

func TestRouter_shortenText(t *testing.T) {
	// // Тестируется ЗАВИСИМОСТИ!
	// // Проверяем, как хендлер отреагирует на ответы от базы данных.
	// // ОПИСАНИЕ ПОЛЯ ТЕСТИРОВАНИЯ.
	// // Здесь хендлер shortenText отдает в db slug и url, а получает
	// type mockBehavior func(s *mocks.MockShortURLRepo)
	// tests := []struct {
	// 	name string
	// 	r    Router
	// 	args args
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		// Инициализация моков.
	// 		ctrl := gomock.NewController(t)
	// 		// Создаем "ложный" сервис, который "притворяется" реальной бизнес-логикой (интерфейсом ShortURLRepo).
	// 		repo := mocks.NewMockShortURLRepo(ctrl)
	// 		// Передавая параметры (repo, tt.id) мы указываем, что
	// 		// 🕖 ожидаем получить вызов методов сервиса repo, а в качестве аргумента передадим id
	// 		tt.mockBehavior
	// 	})
	// }
}
