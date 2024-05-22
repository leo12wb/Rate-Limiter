package limiter

import (
	"sync"
	"time"
)

type localStore struct {
	counters  map[string]*count
	lastEvict time.Time
	mu        sync.Mutex
}

type count struct {
	value     int64
	validThru time.Time
}

var _ CounterStore = &localStore{}

func newLocalStore() *localStore {
	return &localStore{
		counters:  make(map[string]*count),
		lastEvict: time.Now(),
	}
}

func (s *localStore) Increment(key string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.evict()

	c, ok := s.counters[key]
	if !ok {
		c = &count{
			value:     0,
			validThru: time.Now().UTC().Add(time.Second).Truncate(time.Second),
		}
		s.counters[key] = c
	}

	c.value += 1

	return c.value, nil
}

func (s *localStore) Set(key string, value int64, timeout time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := count{
		value:     value,
		validThru: time.Now().UTC().Add(timeout).Truncate(time.Second),
	}
	s.counters[key] = &c

	return nil
}

func (s *localStore) evict() {
	window := time.Second

	if time.Since(s.lastEvict) < window {
		return
	}
	now := time.Now()
	s.lastEvict = now

	for k, c := range s.counters {
		if now.After(c.validThru) {
			delete(s.counters, k)
		}
	}
}
