package service

import (
	"errors"

	"github.com/EvgenyGerdt/go-link-shortening.git/internal/repository"
)

type LinkRepository interface {
	Create(id, url string) error
	GetByID(id string) (string, error)
}

type Shortener struct {
	repo LinkRepository
}

func NewShortener(repo LinkRepository) *Shortener {
	return &Shortener{
		repo: repo,
	}
}

func (s *Shortener) ShortenURL(url string) (string, error) {
	for range 5 {
		id, err := generateID(8)

		if err != nil {
			return "", err
		}

		if err := s.repo.Create(id, url); err == nil {
			return id, nil
		}

		if !errors.Is(err, repository.ErrAlreadyExists) {
			return "", err
		}
	}
	return "", errors.New("shorten url already exists")
}

func (s *Shortener) ResolveURL(id string) (string, error) {
	return s.repo.GetByID(id)
}
