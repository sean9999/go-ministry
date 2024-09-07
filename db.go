package main

import (
	"encoding/json"

	"github.com/sean9999/go-oracle"
	"github.com/sean9999/harebrain"
)

func saveNode(n *node) error {
	db := harebrain.NewDatabase()
	db.Open("data")
	return db.Table("peers").Insert(n)
}

func loadAllNodeRecords() ([][]byte, error) {
	db := harebrain.NewDatabase()
	db.Open("data")
	kv, err := db.Table("peers").GetAll()
	if err != nil {
		return nil, err
	}
	rows := make([][]byte, 0, len(kv))
	for _, v := range kv {
		rows = append(rows, v)
	}
	return rows, nil
}

func loadAllRelationshipRecords() ([][]byte, error) {
	db := harebrain.NewDatabase()
	db.Open("data")
	kv, err := db.Table("rels").GetAll()
	if err != nil {
		return nil, err
	}
	rows := make([][]byte, 0, len(kv))
	for _, v := range kv {
		rows = append(rows, v)
	}
	return rows, nil
}

func loadAllRelationships() ([]relationship, error) {
	records, err := loadAllRelationshipRecords()
	if err != nil {
		return nil, err
	}
	rels := make([]relationship, 0, len(records))
	var prel relationship
	for _, rec := range records {
		err = prel.UnmarshalBinary(rec)
		if err != nil {
			return nil, err
		}
		rels = append(rels, prel)
	}
	return rels, nil
}

func loadAllNodes() ([]*node, error) {
	records, err := loadAllNodeRecords()
	if err != nil {
		return nil, err
	}
	peers := make([]*node, len(records))
	var pbuf node
	for i, rec := range records {
		err = pbuf.UnmarshalJSON(rec)
		if err != nil {
			return nil, err
		}
		peers[i] = &pbuf
	}
	return peers, nil
}

func saveRelationshipSkinny(rel *relationship) error {
	db := harebrain.NewDatabase()
	db.Open("data")
	return db.Table("rels").Insert(rel)
}

func relationshipExists(rel *relationship) bool {
	db := harebrain.NewDatabase()
	db.Open("data")
	_, err := db.Table("rels").Get(rel.Hash())
	return err == nil
}

func removeRelationship(rel *relationship) {
	db := harebrain.NewDatabase()
	db.Open("data")
	db.Table("rels").Delete(rel.Hash())
}

func saveRelationship(j1, j2 json.RawMessage) error {
	db := harebrain.NewDatabase()
	db.Open("data")
	var p1 oracle.Peer
	var p2 oracle.Peer
	p1.UnmarshalJSON(j1)
	p2.UnmarshalJSON(j2)
	rel := newRelationship(p1, p2)
	return db.Table("rels").Insert(&rel)
}
