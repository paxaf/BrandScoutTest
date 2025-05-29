package usecase

import (
	"errors"
	"log"
	"strconv"

	"github.com/paxaf/BrandScoutTest/internal/entity"
)

func (uc *usecase) Delete(key string) error {
	_, ok := uc.repo.Get(key)
	if !ok {
		return errors.New("no rows affected")
	}
	uc.repo.Del(key)
	return nil
}

func (uc *usecase) Random() (entity.Quote, bool) {
	val, ok := uc.repo.GetRandom()
	if !ok {
		log.Println("database is empty")
	}
	return val, ok
}

func (uc *usecase) GetAllByAuthor(author string) ([]entity.Quote, bool) {
	return uc.repo.GetAllByAuthor(author)
}

func (uc *usecase) GetAll() []entity.Quote {
	return uc.repo.GetAll()
}

func (uc *usecase) Set(value entity.Quote) {
	key := uc.keyCounter.Add(1)
	uc.repo.Set(strconv.Itoa(int(key)), value)
	log.Println("successeful set value")
}
