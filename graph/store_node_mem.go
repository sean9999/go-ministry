package graph

import "errors"

var _ Store[*Node] = (*nodeMemStore)(nil)

type nodeMemStore struct {
	db map[string]*Node
}

func NewNodeMemStore() *nodeMemStore {
	s := nodeMemStore{
		db: make(map[string]*Node, 128),
	}
	return &s
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
