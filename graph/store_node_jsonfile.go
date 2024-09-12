package graph

import (
	"encoding/json"

	"github.com/sean9999/harebrain"
)

// *nodeHareStore implements Store[*Node]
var _ Store[*Node] = (*nodeHareStore)(nil)

type nodeHareStore struct {
	tbl *harebrain.Table
}

type jsonNode struct {
	*Node
}

func (jn *jsonNode) Hash() string {
	return jn.Node.Hash() + ".json"
}
func (jn *jsonNode) MarshalBinary() ([]byte, error) {
	return jn.Node.MarshalBinary()
}

func NewJSONNodeStore() *nodeHareStore {
	db := hareBase()
	tbl := db.Table(NODES_FOLDER)
	return &nodeHareStore{tbl}
}

func (s *nodeHareStore) AllRecords() (map[string][]byte, error) {
	orecs, err := s.tbl.GetAll()
	if err != nil {
		return nil, err
	}
	m := make(map[string][]byte, len(orecs))
	for longkey, record := range orecs {
		shortkey := withoutJsonExt(longkey)
		m[shortkey] = record
	}
	return m, nil
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
	raw, err := s.tbl.Get(k + ".json")
	if err != nil {
		return nil, err
	}
	n := new(Node)
	err = n.UnmarshalJSON(raw)
	return n, err
}

func (s *nodeHareStore) Save(v *Node) error {
	jv := &jsonNode{v}
	return s.tbl.Insert(jv)
}

func (s *nodeHareStore) Has(k key) bool {
	_, err := s.tbl.Get(k + ".json")
	return err != nil
}

func (s *nodeHareStore) Delete(k key) error {
	return s.tbl.Delete(k + ".json")
}
