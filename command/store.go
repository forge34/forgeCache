package command

import "sync"

type MapValue struct {
	Value     any
	ExpiresAt int
}

type Store struct {
	mu    sync.Mutex
	Store map[string]*MapValue
}

func NewStore() *Store {
	return &Store{
		Store: make(map[string]*MapValue),
	}
}

func (s *Store) Set(key string, value *MapValue) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Store[key] = value
}

func (s *Store) Get(key string) *MapValue {
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.Store[key]

	if !ok {
		return nil
	}

	return val
}
