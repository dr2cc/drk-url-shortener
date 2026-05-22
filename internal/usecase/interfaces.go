package usecase

// 📍(сюда вернуться) Stabilization Stage (Production-Ready MVP)
// 5️⃣ Полноценная обработка контекста (context.Context)
type (
	// Интерфейсы самих Use Cases (Сервисов) —
	// контракты, которые нужны транспорту (controller / хэндлеры).
	ShortURL interface {
		// Функцонал:
		// Форматирование ID в результирующую строку
		FormatShortURL(baseURL string, urlID string) string
		// Получаем ID, сохраняем в db (в дальнейшем еще и мапим URL из запроса в структуру, к примеру link.ExpandedURL)
		Shorten(url string) (string, error)
		// Находит в хранилище полный URL-адрес по указанному идентификатору.
		// (в дальнейшем еще и возвращает заполненную структуру link.ExpandedURL)
		GetOriginal(slug string) (string, error)
	}

	// Интерфейсы репозиториев (Repositories) —
	// контракты, которые нужны юзкейсам для работы с БД (в evrone реализует их в пакете usecase/repo).
	// То, что юзкейс требует от базы данных
	ShortURLRepo interface {
		Save(slug, url string) error
		Get(slug string) (string, error)
		// 	SaveURL(alias string, url string) error
		// 	GetURL(alias string) (string, error)
	}
)
