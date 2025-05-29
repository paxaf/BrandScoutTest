package storage

import (
	"log"

	"github.com/paxaf/BrandScoutTest/internal/entity"
)

type Engine struct {
	partitions []*HashTable
}

func NewEngine() (*Engine, error) {
	engine := &Engine{
		partitions: make([]*HashTable, 1),
	}
	engine.partitions[0] = NewHashTable()
	return engine, nil
}

func (e *Engine) Set(key string, value entity.Quote) {
	e.partitions[0].Set(key, value)
	log.Println("succeseful set query")
}

func (e *Engine) Get(key string) (entity.Quote, bool) {
	value, found := e.partitions[0].Get(key)
	log.Println("succesefull get query")
	return value, found
}

func (e *Engine) Del(key string) {
	e.partitions[0].Del(key)
	log.Println("succesefull delete query")
}
