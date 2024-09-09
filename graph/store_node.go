package graph

import (
	"encoding/json"
	"errors"

	"github.com/sean9999/harebrain"
)

var _ Store[*Node] = (*nodeHareStore)(nil)
var _ Store[*Node] = (*nodeMemStore)(nil)

type nodeHareStore struct {
	tbl *harebrain.Table
}

type nodeMemStore struct {
	db map[string]*Node
}

func (nm *nodeMemStore) All() []*Node {
	arr := make([]*Node, 0, len(nm.db))
	for _, node := range nm.db {
		arr = append(arr, node)
	}
	return arr
}

func (nm *nodeMemStore) AllRecords() (map[string][]byte, error) {
	m := make(map[string][]byte, len(nm.db))
	for k, v := range nm.db {
		b, _ := v.MarshalBinary()
		m[k] = b
	}
	return m, nil
}

func (nm *nodeMemStore) Get(key string) (*Node, error) {
	record, exists := nm.db[key]
	if !exists {
		return nil, errors.New("record not found")
	}
	return record, nil
}

func (nm *nodeMemStore) Save(n *Node) error {
	nm.db[n.Hash()] = n
	return nil
}

func (nm *nodeMemStore) Has(key string) bool {
	_, exists := nm.db[key]
	return exists
}

func (nm *nodeMemStore) Delete(key string) error {
	delete(nm.db, key)
	return nil
}

func NewJSONNodeStore() *nodeHareStore {
	db := hareBase()
	tbl := db.Table(NODES_FOLDER)
	return &nodeHareStore{tbl}
}

func (s *nodeHareStore) AllRecords() (map[string][]byte, error) {
	return s.tbl.GetAll()
}

func (s *nodeHareStore) All() []*Node {
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

func (s *nodeHareStore) Get(k key) (*Node, error) {
	raw, err := s.tbl.Get(k)
	if err != nil {
		return nil, err
	}
	n := new(Node)
	err = n.UnmarshalJSON(raw)
	return n, err
}

func (s *nodeHareStore) Save(v *Node) error {
	return s.tbl.Insert(v)
}

func (s *nodeHareStore) Has(k key) bool {
	_, err := s.tbl.Get(k)
	return err != nil
}

func (s *nodeHareStore) Delete(k key) error {
	return s.tbl.Delete(k)
}
