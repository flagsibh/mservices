package server

import (
	"log"
	"net/http"
	"time"

	gh "github.com/gorilla/handlers"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")

// New server
func New(mux http.Handler, l *log.Logger) *http.Server {
	env.Parse()

	l.Printf("Starting server on '%s'", *bindAddress)

	// CORS
	// This could also have been implemented with a middleware:
	//
	// func accessControlMiddleware(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		w.Header().Set("Access-Control-Allow-Origin", "*")
	// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
	// 		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
	//
	// 		if r.Method == "OPTIONS" {
	// 			return
	// 		}
	//
	// 		next.ServeHTTP(w, r)
	// 	})
	// }

	handler := gh.CORS(
		gh.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Language", "Content-Type", "X-Requested-With", "Origin"}),
		gh.AllowedMethods([]string{"HEAD", "GET", "POST", "OPTIONS", "DELETE", "PUT"}),
		gh.AllowedOrigins([]string{"*"}),
	)

	srv := &http.Server{
		Addr:         *bindAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler(mux),
	}
	return srv
}
