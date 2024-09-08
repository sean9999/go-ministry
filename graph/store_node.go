package graph

import (
	"encoding/json"

	"github.com/sean9999/harebrain"
)

type nodeStore struct {
	tbl *harebrain.Table
}

func NewNodeStore() *nodeStore {
	db := storeBase()
	tbl := db.Table(NODES_FOLDER)
	return &nodeStore{tbl}
}

func (s *nodeStore) AllRecords() (map[string][]byte, error) {
	return s.tbl.GetAll()
}

func (s *nodeStore) All() []*Node {
	m, err := s.tbl.GetAll()
	if err != nil {
		panic(err)
	}
	nodes := make([]*Node, 0, len(m))
	n := new(Node)
	for _, v := range m {
		json.Unmarshal(v, n)
		nodes = append(nodes, n)
	}
	return nodes
}

func (s *nodeStore) Get(k key) (*Node, error) {
	raw, err := s.tbl.Get(k)
	if err != nil {
		return nil, err
	}
	n := new(Node)
	err = n.UnmarshalJSON(raw)
	return n, err
}

func (s *nodeStore) Save(v *Node) error {
	return s.tbl.Insert(v)
}

func (s *nodeStore) Has(k key) bool {
	_, err := s.tbl.Get(k)
	return err != nil
}

func (s *nodeStore) Delete(k key) error {
	return s.tbl.Delete(k)
}
