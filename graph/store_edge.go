package graph

import (
	"errors"

	"github.com/sean9999/harebrain"
)

var _ Store[*Edge] = (*edgeHareStore)(nil)
var _ Store[*Edge] = (*edgeMemStore)(nil)

type edgeHareStore struct {
	tbl *harebrain.Table
}

func NewJSONEdgeStore() *edgeHareStore {
	db := hareBase()
	tbl := db.Table(EDGES_FOLDER)
	return &edgeHareStore{tbl}
}

func (s *edgeHareStore) Get(k key) (*Edge, error) {
	raw, err := s.tbl.Get(k)
	if err != nil {
		return nil, err
	}
	e := new(Edge)
	err = e.UnmarshalBinary(raw)
	return e, err
}

func (s *edgeHareStore) Save(e *Edge) error {
	return s.tbl.Insert(e)
}

func (s *edgeHareStore) Has(k key) bool {
	_, err := s.tbl.Get(k)
	return err != nil
}

func (s *edgeHareStore) Delete(k key) error {
	return s.tbl.Delete(k)
}

func (s *edgeHareStore) AllRecords() (map[string][]byte, error) {
	return s.tbl.GetAll()
}

func (s *edgeHareStore) All() []*Edge {
	m, err := s.tbl.GetAll()
	if err != nil {
		panic(err)
	}
	edges := make([]*Edge, 0, len(m))
	e := new(Edge)
	for _, v := range m {
		e.UnmarshalBinary(v)
		edges = append(edges, e)
	}
	return edges
}

type edgeMemStore struct {
	db map[string]*Edge
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
