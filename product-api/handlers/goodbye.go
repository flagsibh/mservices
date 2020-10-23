package handlers

import (
	"log"
	"net/http"
)

// GoodBye type
type GoodBye struct {
	l *log.Logger
}

// NewGoodBye constructor
func NewGoodBye(l *log.Logger) *GoodBye {
	return &GoodBye{l}
}

func (g *GoodBye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Byee"))
}
