package main

import "sync"

type Store struct {
	mu        sync.Mutex
	resources []SumberDaya
	nextResID int
}

func NewStore() *Store {
	return &Store{nextResID: 1}
}

func (s *Store) AddResource(r SumberDaya) SumberDaya {
	s.mu.Lock()
	defer s.mu.Unlock()
	r.ID = s.nextResID
	s.nextResID++
	s.resources = append(s.resources, r)
	return r
}

func (s *Store) ListResources() []SumberDaya {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]SumberDaya, len(s.resources))
	copy(out, s.resources)
	return out
}

func (s *Store) GetResource(id int) (SumberDaya, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.findResource(id)
}

func (s *Store) findResource(id int) (SumberDaya, bool) {
	for _, r := range s.resources {
		if r.ID == id {
			return r, true
		}
	}
	return SumberDaya{}, false
}
