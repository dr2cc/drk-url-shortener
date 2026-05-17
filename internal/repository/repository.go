package repository

// URLSaverGetter описывает только то, что нужно хэндлерам
type URLSaverGetter interface {
	SaveURL(alias string, url string) error
	GetURL(alias string) (string, error)
}

type Repository struct {
	URLSaverGetter
}

func New(repo map[string]string) *Repository {
	// make(map[string]string)
	return &Repository{
		URLSaverGetter: NewCache(repo),
	}
}
