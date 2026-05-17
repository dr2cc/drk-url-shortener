package service

import "drk-url-shortener/internal/repository"

// URLSaverGetter описывает только то, что нужно хэндлерам
type URLSaverGetter interface {
	SaveURL(alias string, url string) error
	GetURL(alias string) (string, error)
}

type Service struct {
	URLSaverGetter
}

func New(repo *repository.Repository) *Service {
	return &Service{}
}
