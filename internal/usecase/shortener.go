package usecase

import (
	"drk-url-shortener/internal/lib/random"
	"fmt"
)

// Slug (слаг) — это часть URL-адреса, которая идентифицирует конкретную страницу или ресурс в человекочитаемом виде.
const slugLength = 6

// ShortenerUseCase -.
type Shortener struct {
	// Доступ к методам интерфейса происходит строго через имя этого поля (s.repo.Save())
	repo ShortURLRepo
}

// New -.
// Передаем интерфейс ShortURLRepo вместо *указателя на репозиторий
func New(repo ShortURLRepo) *Shortener {
	// 🔙return structs»
	// DI. Суть- встраиваем repo (как зависимость?) в сервис
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
