package repository

type Storage interface {
	Create(url, slug string) error
	Read(slug string) (string, error)
}

type Repo struct {
	Repo map[string]string
}

func New() Repo {
	return Repo{
		Repo: make(map[string]string),
	}
}
