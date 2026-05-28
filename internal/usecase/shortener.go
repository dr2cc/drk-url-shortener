package usecase

import (
	"drk-url-shortener/internal/lib/random"
	"fmt"
	"strings"
)

// Slug (слаг) — это часть URL-адреса, которая идентифицирует конкретную страницу или ресурс в человекочитаемом виде.
const slugLength = 6

// ShortenerUseCase -.
// Здесь нужно использовать именованные поля (композиция или явное делегирование).
// Почему:
// Здесь мы возвращаем структуру Shortener, как обект имплементирующий ShortURL interface.
// Мы обязаны явно указать для Shortener все методы ShortURL interface
//
// Суть repo это DI, встраивание repo как зависимости в сервис.
// Мы не «встраиваете» её в терминах языка (как embedding), а именно передаем (внедряеи) извне как зависимость.
// Поле repo здесь выступает в роли «приёмника» этой зависимости.
type Shortener struct {
	// «Accept interfaces (ShortURLRepo)🔜,
	// Доступ к методам интерфейса будет происходить строго через имя этого поля (s.repo.Save())
	repo ShortURLRepo
}

// New -.
// Передаем интерфейс ShortURLRepo вместо *указателя на репозиторий
func New(repo ShortURLRepo) *Shortener {
	// 🔙return structs»
	return &Shortener{
		repo: repo,
	}
}

// *For TESTING* // метод Shorten это Impure Function (функция с побочными эффектами, "нечистая").
//
//	Внутри скрывается обращение к интерфейсу базы данных и зависимость от генератора случайных чисел.
//
// Результат зависит от того, что сейчас лежит в БД и что вернет генератор.
// Именно для тестирования таких методов вам и нужен gomock.
func (s *Shortener) Shorten(url string) (string, error) {
	slug := random.NewRandomString(slugLength)
	// Вызываем контракт базы данных через интерфейс
	err := s.repo.Save(slug, url)
	if err != nil {
		return "", err
	}
	return slug, nil
}

func (s Shortener) GetOriginal(slug string) (string, error) {
	// Получаем оригинальный URL из репозитория
	url, err := s.repo.Get(slug)
	if err != nil {
		return "", err
	}
	return url, nil
}

// FormatShortURL производит форматирование полученного ID (путем конкатенации с BaseURL из cfg)
// в результирующую строку, возвращаемую запросами POST
// *For TESTING* // метод FormatShortURL  это Pure Function (чистая функция).
// Она не зависит от внешнего состояния (баз данных, сети). Она не имеет скрытых зависимостей.
// При одних и тех же входных данных она всегда возвращает одинаковый результат.
// Методика тестирования: Чистые функции всегда тестируются с помощью обычных табличных тестов без использования моков.
func (s Shortener) FormatShortURL(baseURL string, urlID string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	urlID = strings.TrimPrefix(urlID, "/")
	return fmt.Sprintf("%s/%s", baseURL, urlID)
}
