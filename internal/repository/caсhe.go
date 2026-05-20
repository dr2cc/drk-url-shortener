package repository

type Cache struct {
	db map[string]string
}

func NewCache(repo map[string]string) Cache {
	// «Accept interfaces, 🔙return structs».
	// Возвращай структуры: "Создатель" объекта знает о нем всё,
	// поэтому возвращает конкретный тип (тут Cache struct).
	// Это дает вызывающему коду гибкость — он сам решит, в какой интерфейс «обернуть» результат.
	return Cache{
		db: repo,
	}
}

func (c Cache) SaveURL(alias string, url string) error {
	c.db[alias] = url
	return nil
}
func (c Cache) GetURL(alias string) (string, error) {
	url := c.db[alias]
	return url, nil
}
