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

// Event (событие)- структура хранения данных в файле с адресами.
// Две главные причины, почему UUID делают строкой:
// 1. Ограничения числовых типов.
// В примере "uuid": "1" (строка) и это не случайно. Тип int в Go — 32 или 64 бита.
// Но UUID (Universally Unique Identifier) — это 128-битное число.
// 2. Простота и совместимость.
// Читаемость: Строковый UUID (в формате 8-4-4-4-12) легко читать глазами в файле конфигурации или логах.
// Гибкость: Если завтра мы решим применить генератор UUID (UUIDv4 (только цифры?), а потом ulid или cuid, которые содержат буквы),
// то не придется менять тип данных в коде и ломать базу данных.
type Event struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Cache struct {
	// Мьютекс (простейшее взаимное исключение, семафор со счетчиком, равным 1)
	// добавлять обязательно в структуру описывающую хранение данных в мапе!
	mu sync.RWMutex // Read-Write Mutex разделяет права на чтение и запись, что оптимальнее sync.Mutex
	db map[string]string
	// Новое к iter9
	filePath string
	nextUUID int
}

// 17.07.2026 Как правильно работать с конструктором (и почему он возвращает указатель) решил посмотреть у Жашкевича.
// После этого решил вновь пройти его книгу с 05 главы.
// Продолжить здесь-
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
}

// Проверить у джеминая комментарии к этой функции (07.07.26)
// Внутренний метод восстановления данных из файла в map
func (c *Cache) loadFromFile() error {
	file, err := os.OpenFile(c.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	// По окончании чтения из файла- закрываем его.
	defer file.Close()

	// Что такое сканер не знаю, но это тип данных в который читается (ридер?) файл с адресами.
	scanner := bufio.NewScanner(file)
	// Читаем файл:
	for scanner.Scan() {
		var event Event
		// Заполняем event данными из файла:
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return err
		}

		// Восстанавливаем кэш
		c.db[event.ShortURL] = event.OriginalURL

		// Вычисляем следующий UUID, чтобы продолжить нумерацию.
		// Переменная id получается из event.UUID; ЕСЛИ ошибки нет и id БОЛЬШЕ или РАВЕН nextUUID ТО
		// присваиваем nextUUID значение id+1
		// Т.к. nextUUID уже изначально равен 1, то при последней итерации он станет больше id и цикл прервется.
		if id, err := strconv.Atoi(event.UUID); err == nil && id >= c.nextUUID {
			c.nextUUID = id + 1
		}
	}
	// Получается при правильном чтении из файла возвращаем эту ошибку.
	// Уточнить почему?
	return scanner.Err()
}

func (c *Cache) Save(alias string, url string) error {
	// Запись, в этот момент полностью блокируем остальным горутинам (пытающимся его захватить=использующим этот мьютекс) и запись и чтение.
	c.mu.Lock()
	defer c.mu.Unlock() // Если забыть снять блокировку по окончании работы функции, программа намертво зависнет (Deadlock),
	// как только любая другая горутина попытается обратиться "за этим" мьютексом.

	// Если файл используется, пишем сначала туда (Write-Throug логика)
	if c.filePath != "" {
		// Продолжить тут-
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
	// Чтение , в этот момент можно разрешить другим горутинам одновременно читать данные.
	// Чтение блокируется только в том случае, если кто-то пишет.
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
