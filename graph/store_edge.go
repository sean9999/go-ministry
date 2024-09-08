package graph

import (
	"github.com/sean9999/harebrain"
)

type edgeStore struct {
	tbl *harebrain.Table
}

func NewEdgeStore() *edgeStore {
	db := storeBase()
	tbl := db.Table(EDGES_FOLDER)
	return &edgeStore{tbl}
}

func (s *edgeStore) Get(k key) (*Edge, error) {
	raw, err := s.tbl.Get(k)
	if err != nil {
		return nil, err
	}
	e := new(Edge)
	err = e.UnmarshalBinary(raw)
	return e, err
}

func (s *edgeStore) Save(e *Edge) error {
	return s.tbl.Insert(e)
}

func (s *edgeStore) Has(k key) bool {
	_, err := s.tbl.Get(k)
	return err != nil
}

func (s *edgeStore) Delete(k key) error {
	return s.tbl.Delete(k)
}

func (s *edgeStore) AllRecords() (map[string][]byte, error) {
	return s.tbl.GetAll()
}

func (s *edgeStore) All() []*Edge {
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
