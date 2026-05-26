package repository

func (r Repo) Create(url, slug string) error {
	r.Repo[slug] = url
	return nil
}
func (r Repo) Read(slug string) (string, error) {
	url := r.Repo[slug]
	return url, nil
}
