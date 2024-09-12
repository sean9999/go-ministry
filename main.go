package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "embed"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sean9999/go-ministry/graph"
)

//go:embed src/favicon.ico
var faviconBytes []byte

func main() {

	// Initialize structured logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	// Create a new Chi router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// WebSocket endpoint
	g := graph.NewMemGraph()
	ws_path := os.Getenv("WS_PATH")
	if ws_path == "" {
		ws_path = "ws"
	}
	r.Mount(fmt.Sprintf("/%s", ws_path), g.Broker)

	//	source maps
	sourceMaps := http.FileServer(http.Dir("."))
	r.Mount("/src", sourceMaps)

	//	favicon
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write(faviconBytes)
	})

	//	static assets
	staticAssets := http.FileServer(http.Dir("./dist"))
	r.Handle("/*", staticAssets)

	// //	send a hello after 5 seconds
	// go func() {
	// 	time.Sleep(5 * time.Second)
	// 	msg := NewMessage()
	// 	msg.Subject = "jazz"
	// 	msg.Payload = json.RawMessage(fmt.Sprintf("%q", "all your base are belong to us"))
	// 	mother.Outbox <- msg
	// }()

	// //	send marco after 9 seconds
	// go func() {
	// 	time.Sleep(9 * time.Second)
	// 	msg := NewMessage()
	// 	msg.Subject = "marco"
	// 	msg.Payload = json.RawMessage("1")
	// 	mother.Outbox <- msg
	// }()

	//	pulse node
	// go func() {
	// 	time.Sleep(3001 * time.Millisecond)
	// 	n := g.RandomNode()
	// 	msg := graph.NewMessage()
	// 	msg.Payload = json.RawMessage(fmt.Sprintf("%q", n.Peer.Nickname()))
	// 	msg.Subject = "command/pulseNode"
	// 	g.Broker.Outbox <- msg
	// }()

	//	update node

	marf := func() {
		for range 5 {
			time.Sleep(997 * time.Millisecond)
			e := g.RandomEdge()
			msg := graph.NewMessage()
			msg.From = e.From()
			msg.To = e.To()
			attrs := graph.NodeAttributes{
				"color": "orange",
			}
			msg.SetPayload(attrs)
			msg.Subject = "command/updateNode"
			g.Broker.Outbox <- msg
		}
	}

	//	process incoming [Message]s
	go func() {
		for msg := range g.Broker.Inbox {
			log.Logger.Info().Str("subject", msg.Subject).Str("payload", string(msg.Payload)).Msg("message receive")

			switch msg.Subject {
			case "marco", "polo":
				rejoinder := msg.Reply()
				if msg.Subject == "marco" {
					rejoinder.Subject = "polo"
				} else {
					rejoinder.Subject = "marco"
				}
				var count int
				json.Unmarshal(msg.Payload, &count)
				count++
				rejoinder.Payload = json.RawMessage(fmt.Sprintf("%d", count))
				err := msg.Conn.WriteJSON(rejoinder)
				if err != nil {
					log.Err(err)
				}
			case "please/addNode":
				n := g.AddNode()
				msg := graph.NewMessage()
				msg.Subject = "command/addPeer"
				peerAsJson, err := n.MarshalJSON()
				if err != nil {
					panic(err)
				}
				msg.Payload = peerAsJson
				g.Broker.Outbox <- msg

			case "please/addRelationship":
				var names [2]string
				err := json.Unmarshal(msg.Payload, &names)
				if err != nil {
					panic(err)
				}
				err = g.AddEdge(graph.Edge{names[0], names[1]})
				if err != nil {
					panic(err)
				}
				msg.Subject = "command/addRelationship"
				g.Broker.Outbox <- msg

			case "hello":
				log.Info().Str("subject", msg.Subject).Str("uuid", msg.ID.String()).Msgf("%v", msg.Payload)
				msg2 := msg.Reply()
				msg2.Subject = "goodbye"
				g.Broker.Outbox <- msg2

			case "hello/imAwake":
				graph.AddABunchOfNodes(g)
				time.Sleep(1 * time.Second)
				graph.AddABunchOfRandomConnections(g)
				graph.DaisyChainConnections(g)
				startEdge := g.RandomEdge()
				n1, err := g.Store.Nodes.Get(startEdge.From())
				if err != nil {
					panic(err)
				}
				n2, err := g.Store.Nodes.Get(startEdge.To())
				if err != nil {
					panic(err)
				}
				graph.Infectify(g, n1, n2)
				break

			case "hello/imAwake_x":

				//	load graph entities
				records, err := g.Store.Nodes.AllRecords()
				if err != nil {
					panic(err)
				}
				for _, record := range records {
					msg := graph.NewMessage()
					msg.Payload = record
					msg.Subject = "command/addPeer"
					g.Broker.Outbox <- msg
				}

				records, err = g.Store.Edges.AllRecords()
				if err != nil {
					panic(err)
				}
				for _, record := range records {
					msg := graph.NewMessage()
					msg.Payload = record
					msg.Subject = "command/addRelationship"
					g.Broker.Outbox <- msg
				}
				go marf()
				// err = g.Store.Zip("cool.zip")
				// fmt.Println("ZIP", err)

			default:
				log.Info().Str("subject", msg.Subject).Msg("default case")
			}

		}
	}()

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8282"
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Info().Str("addr", addr).Msg("Starting server")
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	if err = srv.ListenAndServe(); err != nil {
		log.Error().Msgf("%v", err)
	}

}
