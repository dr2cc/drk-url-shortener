package main

// func TestShortenText(t *testing.T) {
// 	// 🔸Arrange
// 	// Только инициализируем хранилище (так как оно глобальное для нашего main)
// 	repo = make(map[string]string)

// 	// Создаем тело запроса (body это всегда "поток"- Reader, то что "можно читать")
// 	url := "https://google.com"
// 	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(url))

// 	// Создаем ResponseRecorder (замена реальному http.ResponseWriter)
// 	w := httptest.NewRecorder()

// 	// Вызываем хендлер
// 	ShortenText(w, r)

// 	// Проверяем статус код
// 	if w.Code != http.StatusCreated {
// 		t.Errorf("expected status 201, got %d", w.Code)
// 	}

// 	// Проверяем, что в repo что-то появилось
// 	if len(repo) == 0 {
// 		t.Error("expected URL to be saved in repo")
// 	}
// }

// func TestExpand(t *testing.T) {
// 	repo = map[string]string{"test-id": "https://google.com"}

// 	// При GET-запросе ничего не отправляем в теле, передаем nil
// 	r := httptest.NewRequest(http.MethodGet, "/test-id", nil)
// 	// Эмулируем параметры пути для Go 1.22+,
// 	// chi с v5.0.12+ вроде тоже такую работу поддерживает
// 	r.SetPathValue("id", "test-id")

// 	w := httptest.NewRecorder()
// 	Expand(w, r)

// 	if w.Code != http.StatusTemporaryRedirect {
// 		t.Errorf("expected 307, got %d", w.Code)
// 	}

// 	location := w.Header().Get("Location")
// 	if location != "https://google.com" {
// 		t.Errorf("expected redirect to google, got %s", location)
// 	}
// }
