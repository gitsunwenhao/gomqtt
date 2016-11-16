package service

import "stathat.com/c/consistent"

type Hash struct {
	consistent *consistent.Consistent
}

func (h *Hash) Add(keyip string) {
	h.consistent.Add(keyip)
}

func (h *Hash) Get(opid string) (string, error) {
	return h.consistent.Get(opid)
}

func (h *Hash) Remove(keyip string) {
	h.consistent.Remove(keyip)
}

func NewHash() *Hash {
	consistent := consistent.New()
	h := &Hash{consistent: consistent}
	return h
}
