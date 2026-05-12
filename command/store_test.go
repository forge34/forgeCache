package command

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewStore(t *testing.T) {
	s := NewStore()
	if s.data == nil {
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
	retrieved, ok := s.Get(key)

	if !ok {
		t.Fatalf("expected to find key %s, but got false", key)
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
	v, ok := s.Get("ghost")
	if ok {
		t.Errorf("expected ok=false for non-existent key, got true")
	}
	if v != nil {
		t.Errorf("expected nil pointer for non-existent key, got %v", v)
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
			s.Get("key")
		}()
	}

	wg.Wait()
}

func TestOverwrite(t *testing.T) {
	s := NewStore()
	key := "updateMe"

	s.Set(key, &MapValue{Value: "first"})
	s.Set(key, &MapValue{Value: "second"})

	retrieved, _ := s.Get(key)
	if retrieved.Value != "second" {
		t.Errorf("expected value 'second' after overwrite, got %v", retrieved.Value)
	}
}

func TestDelete(t *testing.T) {
	s := NewStore()
	key := "deleteMe"
	s.Set(key, &MapValue{Value: 1})

	deleted := s.Delete(key)
	if !deleted {
		t.Error("expected true when deleting existing key, got false")
	}

	_, ok := s.Get(key)
	if ok {
		t.Error("expected key to not exist after deletion")
	}
}

func TestLen(t *testing.T) {
	s := NewStore()
	s.Set("k1", &MapValue{Value: 1})
	s.Set("k2", &MapValue{Value: 1})
	s.Set("k3", &MapValue{Value: 1})

	if l := s.Len(); l != 3 {
		t.Errorf("Expected length to equal 3, got %d", l)
	}
}

func TestClear(t *testing.T) {
	s := NewStore()
	s.Set("k1", &MapValue{Value: 1})
	s.Set("k2", &MapValue{Value: 1})
	s.Set("k3", &MapValue{Value: 1})

	s.Clear()

	if l := s.Len(); l != 0 {
		t.Errorf("Expected length to equal 0, got %d", l)
	}
}

func TestDeleteNonExistentKey(t *testing.T) {
	s := NewStore()
	deleted := s.Delete("ghost")
	if deleted {
		t.Error("expected false when deleting non-existent key, got true")
	}
}

func TestLenAfterDelete(t *testing.T) {
	s := NewStore()
	s.Set("k1", &MapValue{Value: 1})
	s.Set("k2", &MapValue{Value: 2})
	s.Delete("k1")
	if l := s.Len(); l != 1 {
		t.Errorf("expected length 1 after delete, got %d", l)
	}
}

func TestClearThenRepopulate(t *testing.T) {
	s := NewStore()
	s.Set("k1", &MapValue{Value: 1})
	s.Set("k2", &MapValue{Value: 2})
	s.Clear()
	s.Set("k3", &MapValue{Value: 3})
	if l := s.Len(); l != 1 {
		t.Errorf("expected length 1 after clear and re-population, got %d", l)
	}
	if v, ok := s.Get("k3"); !ok || v.Value != 3 {
		t.Errorf("expected to find k3 with value 3 after re-population, got %v", v)
	}
}

func TestSetNilValue(t *testing.T) {
	s := NewStore()
	s.Set("k1", nil)
	v, ok := s.Get("k1")
	if !ok {
		t.Fatal("expected key to exist after setting nil value, got ok=false")
	}
	if v != nil {
		t.Errorf("expected nil value, got %v", v)
	}
}

func TestConcurrentMixedOps(t *testing.T) {
	s := NewStore()
	iterations := 1000
	var wg sync.WaitGroup

	wg.Add(iterations)
	for i := range iterations {
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", n)
			s.Set(key, &MapValue{Value: n})
		}(i)
	}

	wg.Add(iterations)
	for i := range iterations {
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", n)
			s.Get(key)
		}(i)
	}

	wg.Add(iterations)
	for i := range iterations {
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", n)
			s.Delete(key)
		}(i)
	}

	wg.Add(10)
	for range 10 {
		go func() {
			defer wg.Done()
			s.Clear()
		}()
	}

	wg.Wait()
}
