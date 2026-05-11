package command

import (
	"sync"
	"testing"
)

func TestNewStore(t *testing.T) {
	s := NewStore()
	if s.Store == nil {
		t.Fatal("expected Store map to be initialized, got nil")
	}
}

func TestSetAndGet(t *testing.T) {
	s := NewStore()
	key := "myKey"
	val := &MapValue{
		Value:     "hello",
		ExpiresAt: 1700000000,
	}

	s.Set(key, val)
	retrieved := s.Get(key)

	if retrieved == nil {
		t.Fatalf("expected to find key %s, but got nil", key)
	}

	if retrieved.Value != val.Value {
		t.Errorf("expected value %v, got %v", val.Value, retrieved.Value)
	}

	if retrieved.ExpiresAt != val.ExpiresAt {
		t.Errorf("expected ExpiresAt %d, got %d", val.ExpiresAt, retrieved.ExpiresAt)
	}
}

func TestGetNonExistentKey(t *testing.T) {
	s := NewStore()
	retrieved := s.Get("ghost")

	if retrieved != nil {
		t.Errorf("expected nil for non-existent key, got %v", retrieved)
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := NewStore()
	iterations := 1000
	var wg sync.WaitGroup

	wg.Add(iterations)
	for i := range iterations {
		go func(n int) {
			defer wg.Done()
			s.Set("key", &MapValue{Value: n})
		}(i)
	}

	wg.Add(iterations)
	for range iterations {
		go func() {
			defer wg.Done()
			_ = s.Get("key")
		}()
	}

	wg.Wait()
}

func TestOverwrite(t *testing.T) {
	s := NewStore()
	key := "updateMe"

	s.Set(key, &MapValue{Value: "first"})
	s.Set(key, &MapValue{Value: "second"})

	retrieved := s.Get(key)
	if retrieved.Value != "second" {
		t.Errorf("expected value 'second' after overwrite, got %v", retrieved.Value)
	}
}
