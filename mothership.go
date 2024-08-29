package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

var addr = flag.String("addr", "localhost:8080", "http service address")

var upg = websocket.Upgrader{}

// our main http.Handler, mounted to "/ws" probably
func (m *MotherShip) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Println("new websocket connection")

	log.Info().Msg("asdf")

	m.Logger.Info().Msg("cool")

	defer func() {
		fmt.Println("closing websocket connection")
	}()

	conn, _ := upg.Upgrade(w, r, nil)
	//defer conn.Close()
	m.Connections[conn] = true
	//defer delete(m.Connections, conn)

	//	receive
	var msg Message
	for {
		err := conn.ReadJSON(&msg)
		m.Logger.Println("receiving", msg)
		fmt.Println("recerereve", msg)

		if err != nil {
			fmt.Println("error reading websocket conn", err)
			break
		}

		msg.Conn = conn
		m.Inbox <- msg
	}

}
