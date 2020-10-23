package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Hello estructure
type Hello struct {
	l *log.Logger
}

// NewHello constructor
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	d, err := ioutil.ReadAll(r.Body)
	if err == nil {
		h.l.Printf("Data: %s", d)
		fmt.Fprintf(rw, "Data: %s", d)
	} else {
		http.Error(rw, "Ooops", http.StatusBadRequest)
	}
}
