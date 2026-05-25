package usecase

import (
	"drk-url-shortener/internal/lib/random"
	"fmt"
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
func (s Shortener) FormatShortURL(baseURL string, urlID string) string {
	return fmt.Sprintf("%s/%s", baseURL, urlID)
}
