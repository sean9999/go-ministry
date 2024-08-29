package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize structured logging
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Create a new Chi router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// WebSocket endpoint
	mother := NewMotherShip()
	r.Mount("/ws", mother)

	//	source maps
	sourceMaps := http.FileServer(http.Dir("."))
	r.Mount("/src", sourceMaps)

	//	www
	staticAssets := http.FileServer(http.Dir("./dist"))
	r.Handle("/*", staticAssets)

	//	send outgoing [Message]s over the correct websocket connections
	go func() {
		for msg := range mother.Outbox {
			log.Logger.Println("outbox", msg)
			if msg.Conn != nil {
				msg.Conn.WriteJSON(msg)
			} else {
				log.Logger.Info().Msgf("this message had no websocket connection: %v", msg)
			}
		}
	}()

	//	process incoming [Message]s
	go func() {
		for msg := range mother.Inbox {
			log.Logger.Info().Msgf("%v", msg)

			switch msg.Subject {
			case "marco":
				newid, _ := uuid.NewV7()
				reply := Message{
					ThreadID: msg.ID,
					ID:       newid,
					Subject:  "polo",
					Conn:     msg.Conn,
				}
				err := msg.Conn.WriteJSON(reply)
				if err != nil {
					log.Err(err)
				}
			case "hello":
				log.Info().Str("hello", msg.Subject).Interface("msg", msg)
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
	addr := fmt.Sprintf(":%s", port)
	log.Info().Str("addr", addr).Msg("Starting server")
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Error().Msgf("%v", err)
	}

}
