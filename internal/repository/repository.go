package repository

import "drk-url-shortener/internal/usecase"

// Переход к паттерну «Кэширующий репозиторий» (Cache-Aside / Read-Through).
// https://share.google/aimode/lZjEqcVN02jB29stb

// Repository — это классический паттерн «Адаптер» или «Прокси».
// Встраивание интерфейса usecase.Repository внутри Repository используется именно для того, чтобы:
// - Автоматически пробросить методы Save и Get от Cache наружу.
// - Избавить себя от написания шаблонного кода (Go boilerplate) — не нужно вручную описывать методы Save и Get для структуры Repository.
// - Сохранить за структурой Repository возможность выступать в роли usecase.Repository.
// В данном слое (инфраструктурном репозитории) такое встраивание оправдано,
// так как задача этого слоя — просто предоставить доступ к данным,
// и нам выгодно автоматически делегировать вызовы нижележащему кэшу или базе данных.
type Repository struct {
	// usecase.Repository описывает только то, что нужно сервису от репозитория (сохранение и нахождение)
	usecase.Repository
	// 🧾 В Go внедрение интерфейса без имени (встраивание или embedding) создает поле,
	// имя которого совпадает с именем самого типа (в данном случае Repository).
}

func New() *Repository {
	// 🔙return structs».
	return &Repository{
		Repository: NewCache(),
	}
}
