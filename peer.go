package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/sean9999/go-oracle"
)

type node struct {
	Peer  oracle.Peer    `json:"peer"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

func (n node) Hash() string {
	return fmt.Sprintf(n.Peer.Nickname(), ".json")
}

func (n *node) MarshalJSON() ([]byte, error) {
	m := make(map[string]any, len(n.Attrs)+2)
	for k, v := range n.Peer.AsMap() {
		//	pub: hex of pubkey
		//	nick: nickname
		m[k] = v
	}
	for k, v := range n.Attrs {
		m[k] = v
	}
	return json.Marshal(m)
}

func (n *node) UnmarshalJSON(b []byte) error {
	var m map[string]any
	err := json.Unmarshal(b, &m)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal node: %w", err)
	}
	pubkey, exists := m["pub"]
	if !exists {
		return errors.New("pub key not found")
	}
	p, err := oracle.PeerFromHex([]byte(pubkey.(string)))
	if err != nil {
		return fmt.Errorf("couldn't hydrate node from hex %q: %w", m["pub"], err)
	}
	n.Peer = p
	delete(m, "pub")
	delete(m, "nick")
	n.Attrs = m
	return nil
}

func (n *node) MarshalBinary() ([]byte, error) {
	return n.MarshalJSON()
}

func (p *node) UnmarshalBinary(b []byte) error {
	return p.UnmarshalJSON(b)
}

func newNode(randy io.Reader) *node {
	p := oracle.New(randy).AsPeer()
	n := node{
		Peer:  p,
		Attrs: map[string]any{},
	}
	return &n
}
