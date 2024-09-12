package graph

import (
	"errors"
)

var _ Store[*Edge] = (*edgeMemStore)(nil)

type edgeMemStore struct {
	db map[string]*Edge
}

func NewEdgeMemStore() *edgeMemStore {
	s := edgeMemStore{
		db: make(map[string]*Edge, 512),
	}
	return &s
}

func (nm *edgeMemStore) All() []*Edge {
	arr := make([]*Edge, 0, len(nm.db))
	for _, node := range nm.db {
		arr = append(arr, node)
	}
	return arr
}

func (nm *edgeMemStore) AllRecords() (map[string][]byte, error) {
	m := make(map[string][]byte, len(nm.db))
	for k, v := range nm.db {
		b, _ := v.MarshalBinary()
		m[k] = b
	}
	return m, nil
}

func (nm *edgeMemStore) Get(key string) (*Edge, error) {
	record, exists := nm.db[key]
	if !exists {
		return nil, errors.New("record not found")
	}
	return record, nil
}

func (nm *edgeMemStore) Save(n *Edge) error {
	nm.db[n.Hash()] = n
	return nil
}

func (nm *edgeMemStore) Has(key string) bool {
	_, exists := nm.db[key]
	return exists
}

func (nm *edgeMemStore) Delete(key string) error {
	delete(nm.db, key)
	return nil
}
