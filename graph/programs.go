package graph

import (
	"time"
)

const TICK = time.Millisecond * 25
const BUNCH = 32

func AddABunchOfNodes(g Graph) {
	for i := range BUNCH {
		time.Sleep(TICK)
		msg := NewMessage()
		n := g.AddNode()
		n.Attrs["num"] = i
		msg.SetPayload(n)
		msg.Subject = "command/addNode"
		g.Broker.Outbox <- msg
	}
}

func OneKing(g Graph) {
	allNodes := g.Store.Nodes.All()
	king := allNodes[0]
	for _, thisNode := range allNodes[1:] {
		e1 := Edge{
			king.Hash(),
			thisNode.Hash(),
		}
		e2 := Edge{
			thisNode.Hash(),
			king.Hash(),
		}
		g.AddEdge(e1)
		msg := NewMessage()
		msg.SetPayload(e1.RawJson())
		msg.Subject = "command/addEdge"
		g.Broker.Outbox <- msg
		g.AddEdge(e2)
		msg = NewMessage()
		msg.SetPayload(e2.RawJson())
		msg.Subject = "command/addEdge"
		g.Broker.Outbox <- msg
	}

}

func AddABunchOfRandomConnections(g Graph) {
	for range BUNCH {
		time.Sleep(TICK)
		n1 := g.RandomNode()
		n2 := g.RandomNode()
		e := Edge{n1.Hash(), n2.Hash()}
		g.AddEdge(e)

		msg := NewMessage()
		msg.SetPayload(e.RawJson())
		msg.Subject = "command/addEdge"
		g.Broker.Outbox <- msg
	}
}

func DaisyChainConnections(g Graph) {

	//	each node is connected to the previous

	nodes := g.Store.Nodes.All()

	for i := range nodes[1:] {
		time.Sleep(TICK)
		j := i + 1
		n1 := nodes[i].Hash()
		n2 := nodes[j].Hash()
		e := Edge{n1, n2}
		g.AddEdge(e)
		msg := NewMessage()
		msg.SetPayload(e.RawJson())
		msg.Subject = "command/addEdge"
		g.Broker.Outbox <- msg
	}

	//	the last one is connected to the first
	e := Edge{
		nodes[len(nodes)-1].Hash(),
		nodes[0].Hash(),
	}
	g.AddEdge(e)
	msg := NewMessage()
	msg.SetPayload(e.RawJson())
	msg.Subject = "command/addEdge"
	g.Broker.Outbox <- msg

}

func Infectify(g Graph, fromNode, toNode *Node) {

	if fromNode == nil {
		return
	}
	if toNode == nil {
		return
	}

	if toNode.Attrs["color"] == "orange" {
		return
	}

	time.Sleep(TICK)
	attrs := NodeAttributes{
		"color":          "orange",
		"_originalColor": "orange",
	}
	g.UpdateNode(toNode.Hash(), attrs)
	msg := NewMessage()

	msg.From = fromNode.Hash()
	msg.To = toNode.Hash()

	msg.Subject = "command/passItOn"

	msg.SetPayload(attrs)
	g.Broker.Outbox <- msg

	outgoingEdges, err := g.OutgoingEdges(toNode.Hash())
	if err != nil {
		panic(err)
	}

	time.Sleep(TICK * 25)

	//fmt.Println("outgoing edges", outgoingEdges)
	//fmt.Println("hash", toNode.Hash())

	for _, e := range outgoingEdges {
		targetNode, err := g.Store.Nodes.Get(e.To())
		if err != nil {
			panic(err)
		}
		go Infectify(g, toNode, targetNode)
	}

}
