package command

import (
	"strconv"
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

func (s *Store) IncrDecr(key string, decrease bool) (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.data[key]

	var current int64

	if !ok {
		current = 0
	} else {
		switch val := v.Value.(type) {
		case int:
			current = int64(val)

		case int64:
			current = val

		case string:
			n, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return 0, false
			}

			current = n

		default:
			return 0, false
		}
	}

	if decrease {
		current--
	} else {
		current++
	}

	if ok {
		v.Value = current
	} else {
		s.data[key] = &MapValue{
			Value: current,
		}
	}

	return current, true
}

func (s *Store) UpdateExpiry(key string, t time.Time) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.data[key]
	if !ok || v == nil {
		return false
	}

	v.ExpiresAt = t.UnixNano()
	return true
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
