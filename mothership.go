package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// MotherShip contains pointers to connections, handles websockets, and reconciles truth
type MotherShip struct {
	Connections map[*websocket.Conn]bool
	Logger      zerolog.Logger
	Inbox       chan Message
	Outbox      chan Message
}

// a constructor, just in case
func NewMotherShip() *MotherShip {
	ms := MotherShip{
		Connections: map[*websocket.Conn]bool{},
	}
	return &ms
}

var addr = flag.String("addr", "localhost:8080", "http service address")

var upg = websocket.Upgrader{}

// our main http.Handler, mounted to "/ws" probably
func (m *MotherShip) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	conn, _ := upg.Upgrade(w, r, nil)
	defer conn.Close()
	m.Connections[conn] = true
	defer delete(m.Connections, conn)

	//	receive
	var msg Message
	for {
		err := conn.ReadJSON(msg)
		if err != nil {
			m.Logger.Println("error reading websocket connection:", err)
			break
		}
		msg.Conn = conn
		m.Inbox <- msg
	}

}
