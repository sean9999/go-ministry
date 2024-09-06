package main

import (
	"encoding/json"
	"fmt"

	"github.com/sean9999/go-oracle"
)

type relationship [2]string

func (rel *relationship) Hash() string {
	return fmt.Sprintf("%s-to-%s.json", rel[0], rel[1])
}

func (rel *relationship) MarshalBinary() ([]byte, error) {
	return json.Marshal(rel)
}

func (rel *relationship) UnmarshalBinary(p []byte) error {
	return json.Unmarshal(p, rel)
}

func newRelationship(from oracle.Peer, to oracle.Peer) relationship {
	var rel relationship
	rel[0] = from.Nickname()
	rel[1] = to.Nickname()
	return rel
}
