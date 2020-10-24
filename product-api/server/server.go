package server

import (
	"net/http"
	"time"

	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

// New server
func New(mux *http.ServeMux) *http.Server {
	env.Parse()

	srv := &http.Server{
		Addr:         *bindAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	return srv
}
