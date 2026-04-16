package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenText(t *testing.T) {
	// 🔸Arrange
	// Инициализируем хранилище (так как оно глобальное)
	repo = make(map[string]string)

	// Создаем тело запроса
	url := "https://google.com"
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(url))

	// Создаем ResponseRecorder (замена реальному http.ResponseWriter)
	rr := httptest.NewRecorder()

	// Вызываем хендлер
	ShortenText(rr, req)

	// Проверяем статус код
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}

	// Проверяем, что в repo что-то появилось
	if len(repo) == 0 {
		t.Error("expected URL to be saved in repo")
	}
}

func TestExpand(t *testing.T) {
	repo = map[string]string{"test-id": "https://yandex.ru"}

	req := httptest.NewRequest(http.MethodGet, "/test-id", nil)
	// Эмулируем параметры пути для Go 1.22+
	req.SetPathValue("id", "test-id")

	rr := httptest.NewRecorder()
	Expand(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("expected 307, got %d", rr.Code)
	}

	location := rr.Header().Get("Location")
	if location != "https://yandex.ru" {
		t.Errorf("expected redirect to yandex, got %s", location)
	}
}
