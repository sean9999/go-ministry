package graph

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"sync"

	"github.com/sean9999/go-oracle"
)

type NodeAttributes map[string]any

func (attr NodeAttributes) MarshalJSON() ([]byte, error) {
	m := map[string]any{}
	maps.Copy(m, attr)
	return json.Marshal(m)
}

func (a1 NodeAttributes) Combine(a2 NodeAttributes) NodeAttributes {
	a3 := make(NodeAttributes, len(a1)+len(a2))
	for k, v := range a1 {
		a3[k] = v
	}
	for k, v := range a2 {
		a3[k] = v
	}
	return a3
}

type Node struct {
	*sync.Mutex
	Peer  oracle.Peer    `json:"peer"`
	Attrs NodeAttributes `json:"attrs,omitempty"`
}

func (n *Node) AsMessage() (*Message, error) {
	m := NewMessage()
	j, err := n.MarshalJSON()
	if err != nil {
		return nil, err
	}
	m.Payload = j
	return &m, nil
}

func (n *Node) Update(attrs NodeAttributes) {
	//n.Lock()
	n.Attrs = n.Attrs.Combine(attrs)
	//n.Unlock()
}

func (n *Node) Hash() string {
	return n.Peer.Nickname()
}

func (n *Node) MarshalJSON() ([]byte, error) {
	m := make(map[string]any, len(n.Attrs)+2)
	for k, v := range n.Peer.AsMap() {
		//	pub: hex of pubkey
		//	nick: nickname
		m[k] = v
	}
	for k, v := range n.Attrs {
		m[k] = v
	}
	return json.MarshalIndent(m, "", "\t")
}

func (n *Node) UnmarshalJSON(b []byte) error {
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

func (n *Node) MarshalBinary() ([]byte, error) {
	return n.MarshalJSON()
}

func (p *Node) UnmarshalBinary(b []byte) error {
	return p.UnmarshalJSON(b)
}

func newNode(randy io.Reader) *Node {
	p := oracle.New(randy).AsPeer()
	n := Node{
		Peer:  p,
		Attrs: map[string]any{},
	}
	return &n
}
