package repository

type Cache struct {
	db map[string]string
}

func NewCache() Cache {
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
	url := c.db[alias]
	return url, nil
}
