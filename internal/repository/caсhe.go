package repository

import (
	"bufio"
	"drk-url-shortener/internal/usecase"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

// Event описывает формат строки в JSON-файле
type Event struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Cache struct {
	// Мапы в Go не потокобезопасны при конкурентной записи. Добавим sync.RWMutex
	mu       sync.RWMutex
	db       map[string]string
	filePath string
	nextUUID int
}

func newCache(filePath string) (*Cache, error) {
	c := &Cache{
		db:       make(map[string]string),
		filePath: filePath,
		nextUUID: 1,
	}

	// Если путь к файлу не задан, работаем только в памяти
	if filePath == "" {
		return c, nil
	}

	// Загружаем данные из файла при старте
	if err := c.loadFromFile(); err != nil {
		return nil, fmt.Errorf("failed to load data from file: %w", err)
	}
	// «Accept interfaces, 🔙return structs».
	// Возвращаем структуру (Cache), реализующую интерфейс usecase.Repository
	return c, nil
	// return &Cache{
	// 	db: make(map[string]string),
	// }, nil
}

// Внутренний метод восстановления данных из файла в map
func (c *Cache) loadFromFile() error {
	file, err := os.OpenFile(c.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return err
		}

		// Заполняем кэш
		c.db[event.ShortURL] = event.OriginalURL

		// Вычисляем следующий UUID, чтобы продолжить нумерацию
		if id, err := strconv.Atoi(event.UUID); err == nil && id >= c.nextUUID {
			c.nextUUID = id + 1
		}
	}

	return scanner.Err()
}

func (c *Cache) Save(alias string, url string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если файл используется, пишем сначала туда (Write-Through логика)
	if c.filePath != "" {
		file, err := os.OpenFile(c.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		defer file.Close()

		event := Event{
			UUID:        strconv.Itoa(c.nextUUID),
			ShortURL:    alias,
			OriginalURL: url,
		}

		if err := json.NewEncoder(file).Encode(&event); err != nil {
			return err
		}
		c.nextUUID++
	}

	// Сохраняем в кэш
	c.db[alias] = url
	return nil
}
func (c *Cache) Get(alias string) (string, error) {
	// С появлением сохранения в файл добавляется только это
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Синтаксис url := c.db[alias] не возвращает ошибку в привычном виде error.
	// Вместо этого мапа возвращает второе булево значение (обычно его называют ok),
	// которое показывает, есть ли такой ключ в мапе.
	url, ok := c.db[alias]
	if !ok {
		return "", usecase.ErrStorageNotFound
	}

	return url, nil
}
