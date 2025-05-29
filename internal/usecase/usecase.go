package usecase

import (
	"sync/atomic"

	"github.com/paxaf/BrandScoutTest/internal/entity"
	"github.com/paxaf/BrandScoutTest/internal/repo"
)

type Usecase interface {
	Delete(key string) error
	Random() (entity.Quote, bool)
	GetAllByAuthor(author string)
	GetAll() []entity.Quote
}

type usecase struct {
	repo       repo.Repository
	keyCounter atomic.Int64
}

func New(repo repo.Repository) *usecase {
	return &usecase{
		repo:       repo,
		keyCounter: atomic.Int64{},
	}
}
