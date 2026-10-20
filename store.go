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
