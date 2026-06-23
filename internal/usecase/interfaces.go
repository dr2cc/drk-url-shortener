package usecase

// 📍(сюда вернуться) Stabilization Stage (Production-Ready MVP)
// 5️⃣ context.Context - Полноценная обработка контекста

//go:generate mockgen -source=interfaces.go -destination=mocks/mock_interfaces.go -package=mocks
type (
	// UseCase — контракт интерактора, который нужен транспорту (controller / хэндлеры).
	UseCase interface {
		// Функцонал:
		// Форматирование ID в результирующую строку
		FormatShortURL(baseURL string, urlID string) string
		// Shorten получает полный URL и возвращает сгенерированный короткий код
		Shorten(url string) (string, error)
		// GetOriginal находит в хранилище полный URL-адрес по указанному коду
		// (в дальнейшем еще и возвращает заполненную структуру link.ExpandedURL или как решу назвать)
		GetOriginal(slug string) (string, error)
	}

	// Технический сервис, который занимается только математикой/рандомом
	CodeGenerator interface {
		RandomString() string
	}

	// Интерфейсы репозиториев (Repositories) —
	// контракты, которые нужны юзкейсам для работы с БД (в evrone реализуют их в пакете usecase/repo).
	// То, что юзкейс требует от базы данных
	Repository interface {
		Save(slug, url string) error
		Get(slug string) (string, error)
	}
)
