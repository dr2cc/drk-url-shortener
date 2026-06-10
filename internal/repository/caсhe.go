package repository

import (
	"drk-url-shortener/internal/usecase"
)

type Cache struct {
	db map[string]string
}

func newCache() Cache {
	// «Accept interfaces, 🔙return structs».
	// Возвращаем структуру (Cache), реализующую интерфейс usecase.ShortURLRepo
	return Cache{
		db: make(map[string]string),
	}
}

func (c Cache) Save(alias string, url string) error {
	c.db[alias] = url
	return nil
}
func (c Cache) Get(alias string) (string, error) {
	// Синтаксис url := c.db[alias] не возвращает ошибку в привычном виде error.
	// Вместо этого мапа возвращает второе булево значение (обычно его называют ok),
	// которое показывает, есть ли такой ключ в мапе.
	url, ok := c.db[alias]
	if !ok {
		return "", usecase.ErrStorageNotFound
	}

	return url, nil
}
