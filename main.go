package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
	mother := NewMotherShip()
	ws_path := os.Getenv("WS_PATH")
	if ws_path == "" {
		ws_path = "ws"
	}
	r.Mount(fmt.Sprintf("/%s", ws_path), mother)

	//	source maps
	sourceMaps := http.FileServer(http.Dir("."))
	r.Mount("/src", sourceMaps)

	//	favicon
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		fd, err := os.Open("./src/favicon.ico")
		if err != nil {
			w.WriteHeader(404)
			w.Write([]byte("favicon not found"))
			return
		}
		defer fd.Close()
		io.Copy(w, fd)
	})

	//	static assets
	staticAssets := http.FileServer(http.Dir("./dist"))
	r.Handle("/*", staticAssets)

	//	send a hello after 5 seconds
	go func() {
		time.Sleep(5 * time.Second)
		msg := NewMessage()
		msg.Subject = "jazz"
		msg.Payload = json.RawMessage(fmt.Sprintf("%q", "all your base are belong to us"))
		mother.Outbox <- msg
	}()

	//	send marco after 9 seconds
	go func() {
		time.Sleep(9 * time.Second)
		msg := NewMessage()
		msg.Subject = "marco"
		msg.Payload = json.RawMessage("1")
		mother.Outbox <- msg
	}()

	//	process incoming [Message]s
	go func() {
		for msg := range mother.Inbox {
			log.Logger.Info().Msgf("%v", msg)

			switch msg.Subject {
			case "marco", "polo":
				rejoinder := msg.Reply()
				if msg.Subject == "marco" {
					rejoinder.Subject = "polo"
				} else {
					rejoinder.Subject = "marco"
				}
				count := new(int)
				json.Unmarshal(msg.Payload, count)
				*count++
				rejoinder.Payload = json.RawMessage(fmt.Sprintf("%d", *count))
				err := msg.Conn.WriteJSON(rejoinder)
				if err != nil {
					log.Err(err)
				}
			case "hello":
				log.Info().Str("subject", msg.Subject).Str("uuid", msg.ID.String()).Msgf("%v", msg.Payload)
				msg2 := msg.Reply()
				msg2.Subject = "goodbye"
				mother.Outbox <- msg2
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
