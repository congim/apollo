package common

import "stathat.com/c/consistent"

// Hash hash struct
type Hash struct {
	consistent *consistent.Consistent
}

// Add add new key
func (h *Hash) Add(key string) {
	h.consistent.Add(key)
}

// Get get value by key
func (h *Hash) Get(key string) (string, error) {
	return h.consistent.Get(key)
}

// Remove remove key
func (h *Hash) Remove(key string) {
	h.consistent.Remove(key)
}

// NewHash return new hash struct
func NewHash() *Hash {
	consistent := consistent.New()
	h := &Hash{consistent: consistent}
	return h
}
