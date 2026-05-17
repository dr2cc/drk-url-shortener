package repository

type Cache struct {
	db map[string]string
}

func NewCache(repo map[string]string) Cache {
	// Возвращай структуры!
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
