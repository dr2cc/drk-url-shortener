package handler

import (
	"net/http"
	"testing"
)

func Test_redirect(t *testing.T) {
	// Создаем структуру для ожидаемого результата
	type wantResult struct {
		statusCode int
		// заголовок
		location string
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		repo map[string]string
		want wantResult
		// Какой URL мы запрашиваем в тесте
		requestPath string
		// Значение именованного параметра {id}
		pathValue string
	}{
		{
			name:        "Positive",
			repo:        map[string]string{"go": "https://go.dev"},
			want:        wantResult{statusCode: http.StatusTemporaryRedirect, location: "https://go.dev"},
			requestPath: "/go",
			// Хендлер будет искать этот ключ в репозитории
			pathValue: "go",
		},
		{
			name:        "Error in path",
			repo:        map[string]string{"go": "https://go.dev"},
			want:        wantResult{statusCode: http.StatusBadRequest, location: ""},
			requestPath: "/java",
			pathValue:   "java",
		},
		{
			// В реальной работе такого не случится. Роутер отправит "/" на POST хендлер
			// Но в тесте можно проверить- тут обращение напрямую к хендлеру.
			name:        "Error: Empty path (Root URL)",
			repo:        map[string]string{"go": "https://go.dev"},
			want:        wantResult{statusCode: http.StatusBadRequest, location: ""},
			requestPath: "/", // Запрос на корень
			pathValue:   "",  // Роутер не найдет параметр {id}
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 	got := redirect(tt.repo)

			// 	// для GET-запроса редиректа передавать данные в теле (например, JSON или форму) не нужно.
			// 	req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			// 	// --- ЭМУЛЯЦИЯ РАБОТЫ РОУТЕРА ДЛЯ r.PathValue ---
			// 	// Вшиваем параметр "id" в контекст запроса, как это сделал бы роутер
			// 	req.SetPathValue("id", tt.pathValue)
			// 	// Действует как браузер. Записывает заголовки, статус-код и тело ответа, которые возвращает обработчик.
			// 	w := httptest.NewRecorder()

			// 	// Выполняем запрос через роутер
			// 	got.ServeHTTP(w, req)

			// 	// Проверяем статус-код
			// 	if w.Code != tt.want.statusCode {
			// 		t.Errorf("redirect() status = %v, want %v", w.Code, tt.want.statusCode)
			// 	}

			// 	// Проверяем заголовок Location (куда редиректит)
			// 	if loc := w.Header().Get("Location"); loc != tt.want.location {
			// 		t.Errorf("redirect() location = %v, want %v", loc, tt.want.location)
			// 	}
		})
	}
}
