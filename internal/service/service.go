package service

import "drk-url-shortener/internal/repository"

// URLSaverGetter описывает только то, что нужно от слоя service
// type ShortURL interface
type URLSaverGetter interface {
	// Функцонал:
	// Форматирование ID в результирующую строку
	FormatShortURL(baseURL string, urlID string) string
	// Получаем ID, сохраняем в db (в дальнейшем еще и мапим URL из запроса в структуру, к примеру link.ExpandedURL)
	ShortenURL(url string) (string, error)
	// Находит в хранилище полный URL-адрес по указанному идентификатору.
	// (в дальнейшем еще и возвращает заполненную структуру link.ExpandedURL)
	FindURL(alias string) (string, error)
}

type Service struct {
	URLSaverGetter
}

func New(repo *repository.Repository) *Service {
	return &Service{
		// «Accept interfaces🔜,
		URLSaverGetter: NewShortener(repo),
	}
}
