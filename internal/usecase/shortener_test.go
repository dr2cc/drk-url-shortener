package usecase

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestShortener_FormatShortURL(t *testing.T) {
	// 1. (📦Arrange) Определение структуры аргументов (Input Data)
	// Эти параметры передаются непосредственно в тестируемый метод.
	type args struct {
		baseURL string
		urlID   string
	}
	// 2. (📦Arrange) Определение структуры тест-кейса (Test Case Definition)
	// Здесь описывается «форма» одного теста.
	tests := []struct {
		name string
		s    Shortener
		args args
		want string
	}{
		// 3. (📦Arrange) Слайс тест-кейсов (Test Cases Table)
		// Конкретные наборы данных: позитивные, негативные, граничные сценарии.
		{
			name: "OK",
			args: args{baseURL: "http://localhost:8080", urlID: "abc"},
			want: "http://localhost:8080/abc",
		},
		{
			name: "empty",
			args: args{baseURL: "", urlID: ""},
			want: "/",
		},
		{
			name: "slash in baseURL",
			args: args{baseURL: "http://localhost:8080/", urlID: "abc"},
			want: "http://localhost:8080/abc",
		},
		{
			name: "slash in urlID",
			args: args{baseURL: "http://localhost:8080", urlID: "/abc"},
			want: "http://localhost:8080/abc",
		},
	}
	// 4. (⚙️Инфраструктура Go, НЕ часть логики теста 📦A🎬A🚨A) Итератор или Тестовый цикл (Test Runner)
	for _, tt := range tests {
		// 5. (⚙️Инфраструктура Go, НЕ часть логики теста 📦A🎬A🚨A) Запуск подтеста (Subtest)
		// t.Run создает изолированную область выполнения для каждого набора данных.
		t.Run(tt.name, func(t *testing.T) {
			// 7. 🚨Assert (проверка ожидаемого результата с фактическим).
			// tt.s.FormatShortURL(tt.args.baseURL, tt.args.urlID) - это 6. 🎬Act (вызов тестируемого действия)
			assert.Equal(t, tt.want, tt.s.FormatShortURL(tt.args.baseURL, tt.args.urlID))
		})
	}
}
