package graph

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	mathrand "math/rand"
)

type Graph struct {
	Broker *Broker
	Store  GraphStore
}

func (g Graph) AddNode() *Node {
	n := newNode(rand.Reader)
	g.Store.Nodes.Save(n)
	return n
}

func (g Graph) RandomNode() *Node {
	a := g.Store.Nodes.All()
	mathrand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a[0]
}

func (g Graph) RandomEdge() *Edge {
	a := g.Store.Edges.All()
	mathrand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a[0]
}

func (g Graph) AddEdge(name1, name2 string) error {
	e := Edge{name1, name2}
	return g.Store.Edges.Save(&e)
}

func (g Graph) UpdateNode(id string, attrs NodeAttributes) error {
	n, err := g.Store.Nodes.Get(id)
	if err != nil {
		return err
	}
	n.Attrs.Combine(attrs)
	return nil
}

func (g Graph) SendMessage(msg Message) error {

	if !g.EdgeExists(msg.From, msg.To) {
		return errors.New("edge doesn't exist")
	}

	pp := new(NodeAttributes)
	err := json.Unmarshal(msg.Payload, pp)
	if err != nil {
		return err
	}
	attrs, convertable := msg.GetPayload().(NodeAttributes)
	if !convertable {
		return fmt.Errorf("%v is not convertable to NodeAttributes", string(msg.Payload))
	}
	err = g.UpdateNode(msg.To, attrs)
	if err != nil {
		return err
	}

	msg2 := msg.Reply()
	msg2.Subject = "command/updateNode"
	g.Broker.Outbox <- msg2
	return nil
}

func (g Graph) EdgeExists(n1, n2 string) bool {
	e := Edge{n1, n2}
	return g.Store.Edges.Has(e.Hash())
}

func NewMemGraph() Graph {
	ship := NewMotherShip()
	pers := GraphStore{
		&nodeMemStore{},
		&edgeMemStore{},
	}
	g := Graph{
		ship, pers,
	}
	return g
}

func NewJSONGraph() Graph {
	ship := NewMotherShip()
	pers := GraphStore{
		NewJSONNodeStore(),
		NewJSONEdgeStore(),
	}
	g := Graph{
		ship, pers,
	}
	return g
}
