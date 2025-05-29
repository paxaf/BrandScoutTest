package storage

import (
	"sync"

	"github.com/paxaf/BrandScoutTest/internal/entity"
)

type HashTable struct {
	mutex sync.RWMutex
	data  map[string]entity.Quote
}

func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]entity.Quote),
	}
}

func (h *HashTable) Set(key string, value entity.Quote) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.data[key] = value
}

func (h *HashTable) Del(key string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.data, key)
}

func (h *HashTable) Get(key string) (entity.Quote, bool) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	value, found := h.data[key]
	return value, found
}
