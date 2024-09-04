package main

import (
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// MotherShip brokers websocket connections and channels
type MotherShip struct {
	Connections map[*websocket.Conn]bool
	Logger      zerolog.Logger
	Inbox       chan Message
	Outbox      chan Message
}

// constructor
func NewMotherShip() *MotherShip {

	z := zerolog.New(os.Stdout)
	z.Level(zerolog.DebugLevel)

	ms := MotherShip{
		Connections: map[*websocket.Conn]bool{},
		Logger:      z,
		Inbox:       make(chan Message, 1024),
		Outbox:      make(chan Message, 1024),
	}
	return &ms
}

var upg = websocket.Upgrader{}

// our main http.Handler, mounted to "/ws" probably
func (m *MotherShip) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	m.Logger.Info().Msg("opening websocket connection")
	conn, _ := upg.Upgrade(w, r, nil)
	m.Connections[conn] = true

	defer func() {
		m.Connections[conn] = false
		m.Logger.Info().Msg("closing websocket connection")
		delete(m.Connections, conn)
		conn.Close()
	}()

	//	send outgoing [Message]s
	go func() {
		for msg := range m.Outbox {
			if msg.Conn != nil {
				//	unicast
				msg.Conn.WriteJSON(msg)
			} else {
				//	broadcast
				for c, is := range m.Connections {
					if is {
						c.WriteJSON(msg)
					}
				}
			}
		}
	}()

	//	receive [Message]s over websocket and queue them on Inbox
	var msg Message
	for {
		err := conn.ReadJSON(&msg)
		m.Logger.Info().Str("subject", msg.Subject).Str("uuid", msg.ID.String()).Msg("receive message")
		if err != nil {
			break
		}
		msg.Conn = conn
		m.Inbox <- msg
	}

}
