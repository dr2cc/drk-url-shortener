package repository

import "drk-url-shortener/internal/usecase"

type Repository struct {
	// «Accept interfaces (usecase.ShortURLRepo)🔜, return structs».
	// usecase.ShortURLRepo описывает только то, что нужно сервису от репозитория (сохранение и нахождение)
	usecase.ShortURLRepo
	// 🧾 В Go внедрение интерфейса без имени (встраивание или embedding) создает поле,
	// имя которого совпадает с именем самого типа (в данном случае ShortURLRepo).
}

func New(repo map[string]string) *Repository {
	// 🔙return structs».
	return &Repository{
		ShortURLRepo: NewCache(repo),
	}
}
