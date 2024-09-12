package graph

import (
	"github.com/sean9999/harebrain"
)

var _ Store[*Edge] = (*edgeHareStore)(nil)

// jsonEdge is an Edge that intercepts its Hash() function for use as JSON files
type jsonEdge struct {
	*Edge
}

func (je *jsonEdge) Hash() string {
	return je.Edge.Hash() + ".json"
}
func (je *jsonEdge) MarshalBinary() ([]byte, error) {
	return je.Edge.MarshalBinary()
}

type edgeHareStore struct {
	tbl *harebrain.Table
}

func NewJSONEdgeStore() *edgeHareStore {
	db := hareBase()
	tbl := db.Table(EDGES_FOLDER)
	return &edgeHareStore{tbl}
}

func (s *edgeHareStore) Get(k key) (*Edge, error) {
	raw, err := s.tbl.Get(k + ".json")
	if err != nil {
		return nil, err
	}
	e := new(Edge)
	err = e.UnmarshalBinary(raw)
	return e, err
}

func (s *edgeHareStore) Save(e *Edge) error {
	je := &jsonEdge{e}
	return s.tbl.Insert(je)
}

func (s *edgeHareStore) Has(k key) bool {
	_, err := s.tbl.Get(k + ".json")
	return err != nil
}

func (s *edgeHareStore) Delete(k key) error {
	return s.tbl.Delete(k + ".json")
}

func (s *edgeHareStore) AllRecords() (map[string][]byte, error) {
	orecs, err := s.tbl.GetAll()
	if err != nil {
		return nil, err
	}
	m := make(map[string][]byte, len(orecs))
	for longkey, v := range orecs {
		shortkey := withoutJsonExt(longkey)
		m[shortkey] = v
	}
	return m, nil
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
