package repo

import "github.com/paxaf/BrandScoutTest/internal/entity"

type Repository interface {
	Set(key string, value entity.Quote)
	Del(key string)
	Get(key string) (entity.Quote, bool)
	GetAllByAuthor(author string) ([]entity.Quote, bool)
	GetRandom() (entity.Quote, bool)
	GetAll() []entity.Quote
}
