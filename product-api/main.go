package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/flagsibh/mservices/product-api/handlers"
	"github.com/flagsibh/mservices/product-api/server"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	ph := handlers.NewProducts(l)

	r := mux.NewRouter()
	getr := r.Methods(http.MethodGet).Subrouter()
	getr.HandleFunc("/", ph.GetProducts)

	putr := r.Methods(http.MethodPut).Subrouter()
	putr.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putr.Use(ph.ProductValidation)

	postr := r.Methods(http.MethodPost).Subrouter()
	postr.HandleFunc("/", ph.AddProduct)
	postr.Use(ph.ProductValidation)

	srv := server.New(r, l)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	sc := <-channel
	l.Println("Received terminate, graceful shutdown", sc)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(tc)
}
