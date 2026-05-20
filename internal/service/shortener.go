package service

import (
	"drk-url-shortener/internal/lib/random"
	"drk-url-shortener/internal/repository"
	"fmt"
)

const aliasLength = 6

type Shortener struct {
	repo *repository.Repository
}

func NewShortener(repo *repository.Repository) *Shortener {
	// 🔙return structs».
	// Возвращай структуры: "Создатель" объекта знает о нем всё,
	// поэтому возвращает конкретный тип (тут Cache struct).
	// Это дает вызывающему коду гибкость — он сам решит, в какой интерфейс «обернуть» результат.
	return &Shortener{
		repo: repo,
	}
}

func (s Shortener) ShortenURL(url string) (string, error) {
	alias := random.NewRandomString(aliasLength)
	s.repo.SaveURL(alias, url)
	return alias, nil
}
func (s Shortener) FindURL(alias string) (string, error) {

	url, _ := s.repo.GetURL(alias)
	return url, nil
}

// FormatShortURL производит форматирование полученного ID (путем конкатенации с BaseURL из cfg)
// в результирующую строку, возвращаемую запросами POST
func (s Shortener) FormatShortURL(baseURL string, urlID string) string {
	return fmt.Sprintf("%s/%s", baseURL, urlID)
}
