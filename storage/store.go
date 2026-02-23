package storage

import "sync"

type InMemoryStore struct {
	db map[string]string
	mu sync.RWMutex
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		db: make(map[string]string),
	}
}

func (s *InMemoryStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if val, ok := s.db[key]; ok {
		return val, true
	}
	return "", false
}

func (s *InMemoryStore) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db[key] = value
}

func (s *InMemoryStore) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.db, key)
	return true
}
