package server

import (
	"net/http"
	"time"
)

// New server
func New(mux *http.ServeMux, serviceAddr string) *http.Server {
	srv := &http.Server{
		Addr:         serviceAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	return srv
}
