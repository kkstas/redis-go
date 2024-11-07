package store

import (
	"log"
	"sync"
	"time"
)

type Value struct {
	value      string
	withExpiry bool
	expiry     int
}

type Store struct {
	values map[string]Value
	mu     sync.RWMutex
}

func New() *Store {
	return &Store{values: make(map[string]Value)}

}

func (s *Store) Set(key, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = Value{value: val}
	log.Printf("set: %#v\n", s.values)
}

func (s *Store) SetWithExpiry(key, val string, expiry int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = Value{
		value:      val,
		withExpiry: true,
		expiry:     int(time.Now().UnixMilli()) + expiry,
	}
	log.Printf("setWithExpiry: %#v\n", s.values)
}

func (s *Store) Get(key string) (string, bool) {
	val, found := s.getValue(key)

	if !found {
		return "", false
	}

	if val.withExpiry && val.expiry < int(time.Now().UnixMilli()) {
		s.deleteKey(key)
		return "", false
	}

	return val.value, true
}

func (s *Store) getValue(key string) (Value, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, found := s.values[key]
	log.Printf("getValue: %#v\n", s.values)
	return val, found
}

func (s *Store) deleteKey(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Printf("before deleteKey: %#v\n", s.values)
	delete(s.values, key)
	log.Printf("after deleteKey: %#v\n", s.values)
}
