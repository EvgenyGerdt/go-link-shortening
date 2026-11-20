package repository

import (
	"errors"
	"sync"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type LinkRepository struct {
	mu    sync.RWMutex
	links map[string]string
}

func NewLinkRepository() *LinkRepository {
	return &LinkRepository{
		links: make(map[string]string),
	}
}

func (r *LinkRepository) Create(id, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.links[id]; exists {
		return ErrAlreadyExists
	}

	r.links[id] = url
	return nil
}

func (r *LinkRepository) GetByID(id string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, exists := r.links[id]
	if !exists {
		return "", ErrNotFound
	}

	return url, nil
}
