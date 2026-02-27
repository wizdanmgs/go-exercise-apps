package store

import (
	"sync"
	"url-shortener/internal/model"
)

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]model.URL
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: map[string]model.URL{}}
}

func (m *MemoryStore) Save(url model.URL) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[url.Code] = url
}

func (m *MemoryStore) Get(code string) (model.URL, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	url, ok := m.data[code]
	return url, ok
}

func (m *MemoryStore) Delete(code string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, code)
}
