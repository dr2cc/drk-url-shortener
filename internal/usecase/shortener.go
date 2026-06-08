package usecase

import (
	"errors"
	"fmt"
	"strings"
)

// ShortenerUseCase -.
// Здесь нужно использовать именованные поля (композиция или явное делегирование).
// Почему:
// Здесь мы возвращаем структуру Shortener, как обект имплементирующий интерфесы Repository и CodeGenerator.
// Мы обязаны явно указать для Shortener все методы этих интерфейсов (для доступности как s.generator.RandomString или s.Repo.Save)
//
// Суть (на примере Repo): Repo это DI, встраивание Repo как зависимости в сервис.
// Мы не «встраиваем» её в терминах языка (как embedding), а именно передаем (внедряеи) извне как зависимость.
// Поле Repo здесь выступает в роли «приёмника» этой зависимости.
type Shortener struct {
	// Доступ к методам интерфейса будет происходить строго через имя этого поля (s.Repo.Save())
	Repo      Repository
	generator CodeGenerator
}

// New -.
// Принимай интерфейсы (Repository, CodeGenerator), возвращай структуры (*Shortener)
func New(r Repository, g CodeGenerator) *Shortener {
	// 🔙return structs»
	return &Shortener{
		Repo:      r,
		generator: g,
	}
}

// *For TESTING* // метод Shorten это Impure Function (функция с побочными эффектами, "нечистая").
//
//	Внутри скрывается обращение к интерфейсу базы данных и зависимость от генератора случайных чисел.
//
// Результат зависит от того, что сейчас лежит в БД и что вернет генератор.
// Методика тестирования: именно для тестирования таких методов и нужен gomock.
func (s *Shortener) Shorten(url string) (string, error) {
	// Вызов технического сервиса: генерируем уникальный хэш-код
	slug := s.generator.RandomString()
	// Вызываем контракт базы данных через интерфейс
	err := s.Repo.Save(slug, url)
	if err != nil {
		return "", err
	}
	return slug, nil
}

func (s Shortener) GetOriginal(slug string) (string, error) {
	// Получаем оригинальный URL из репозитория
	url, err := s.Repo.Get(slug)
	if err != nil {
		// Проверяем: если это ошибка "не найдено" из репозитория
		if errors.Is(err, ErrStorageNotFound) {
			// Возвращаем чистую бизнес-ошибку наружу (роутеру)
			return "", ErrCodeNotFound
		}

		// Любую другую техническую ошибку (например, упала сеть к БД)
		// оборачиваем через %w, чтобы не терять контекст для логов
		return "", fmt.Errorf("failed to get url from repository: %w", err)
	}
	return url, nil
}

// *For TESTING* // метод FormatShortURL  это Pure Function (чистая функция).
// Она не зависит от внешнего состояния (баз данных, сети). Она не имеет скрытых зависимостей.
// При одних и тех же входных данных она всегда возвращает одинаковый результат.
// Методика тестирования: чистые функции всегда тестируются с помощью обычных табличных тестов без использования моков.
//
// FormatShortURL производит форматирование полученного urlID (путем конкатенации с BaseURL из cfg)
// в результирующую строку, возвращаемую запросами POST
func (s Shortener) FormatShortURL(baseURL string, urlID string) string {
	baseURL = strings.TrimSuffix(baseURL, "/")
	urlID = strings.TrimPrefix(urlID, "/")
	return fmt.Sprintf("%s/%s", baseURL, urlID)
}
