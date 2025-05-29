package storage

import (
	"log"

	"github.com/paxaf/BrandScoutTest/internal/entity"
)

type Engine struct {
	partition *HashTable
}

func NewEngine() (*Engine, error) {
	engine := &Engine{
		partition: NewHashTable(),
	}
	engine.partition = NewHashTable()
	return engine, nil
}

func (e *Engine) Set(key string, value entity.Quote) {
	e.partition.Set(key, value)
	log.Println("succeseful set query")
}

func (e *Engine) Get(key string) (entity.Quote, bool) {
	value, found := e.partition.Get(key)
	log.Println("succesefull get query")
	return value, found
}

func (e *Engine) Del(key string) {
	e.partition.Del(key)
	log.Println("succesefull delete query")
}

func (e *Engine) GetAllByAuthor(author string) ([]entity.Quote, bool) {
	var res []entity.Quote
	for _, val := range e.partition.data {
		if val.Author == author {
			res = append(res, val)
		}
	}
	if len(res) < 1 {
		return nil, false
	}
	return res, true
}

func (e *Engine) GetRandom() (entity.Quote, bool) {
	for _, val := range e.partition.data {
		return val, true
	}
	return entity.Quote{}, false
}
