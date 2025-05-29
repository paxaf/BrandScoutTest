package usecase

import (
	"github.com/paxaf/BrandScoutTest/internal/repo"
)

type Usecase interface {
}

type usecase struct {
	repo repo.Repository
}

func New(repo repo.Repository) *usecase {
	return &usecase{
		repo: repo,
	}
}
