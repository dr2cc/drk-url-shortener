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
	return &Repository{
		// «Accept interfaces🔜, return structs».
		// Принимай интерфейсы (тут URLSaverGetter interface):
		// Функция должна требовать только то поведение, которое ей реально нужно!
		URLSaverGetter: NewCache(repo),
	}
}
