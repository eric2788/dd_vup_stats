package set

import "sync"

// Create Generic Set

type void struct{}

type Set[K int64 | string | int | int32 | float64] struct {
	mu sync.Mutex
	m  map[K]void
}

func New[K int64 | string | int | int32 | float64]() *Set[K] {
	return &Set[K]{
		m: make(map[K]void),
	}
}

func (s *Set[K]) Add(value K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[value] = void{}
}

func (s *Set[K]) Remove(value K) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, value)
}

func (s *Set[K]) Has(value K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.m[value]
	return ok
}

func (s *Set[K]) Difference(other *Set[K]) *Set[K] {
	s.mu.Lock()
	defer s.mu.Unlock()
	diff := New[K]()
	for k, _ := range s.m {
		if !other.Has(k) {
			diff.Add(k)
		}
	}
	return diff
}

func (s *Set[K]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m = make(map[K]void)
}

// ToArray to array
func (s *Set[K]) ToArray() []K {
	s.mu.Lock()
	defer s.mu.Unlock()
	array := make([]K, 0)
	for k, _ := range s.m {
		array = append(array, k)
	}
	return array
}

func FromArray[K int64 | string | int | int32 | float64](arr []K) *Set[K] {
	s := New[K]()
	for _, v := range arr {
		s.Add(v)
	}
	return s
}
