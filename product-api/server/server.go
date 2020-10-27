package server

import (
	"log"
	"net/http"
	"time"

	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

// New server
func New(mux http.Handler, l *log.Logger) *http.Server {
	env.Parse()

	l.Printf("Starting server on '%s'", *bindAddress)

	srv := &http.Server{
		Addr:         *bindAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	return srv
}
