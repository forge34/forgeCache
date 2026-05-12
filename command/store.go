package command

import (
	"sync"
	"time"
)

type MapValue struct {
	Value     any
	ExpiresAt int64
}

type Store struct {
	mu   sync.RWMutex
	data map[string]*MapValue
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]*MapValue),
	}
}

func (s *Store) Set(key string, value *MapValue) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) Get(key string) (*MapValue, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]

	if !ok {
		return nil, false
	}

	if val.ExpiresAt != 0 && val.ExpiresAt < time.Now().UnixNano() {
		return nil, false
	}

	return val, true
}

func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[key]; !ok {
		return false
	}
	delete(s.data, key)

	return true
}

func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	clear(s.data)
}
