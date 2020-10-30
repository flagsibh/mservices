package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/flagsibh/mservices/product-api/handlers"
	mw "github.com/flagsibh/mservices/product-api/handlers/middleware"
	"github.com/flagsibh/mservices/product-api/server"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	v := mw.NewValidation(l)

	ph := handlers.NewProducts(l)

	r := mux.NewRouter()
	r.Use(mw.ContentTypeMiddleware)

	getr := r.Methods(http.MethodGet).Subrouter()
	getr.HandleFunc("/", ph.GetProducts)
	getr.HandleFunc("/{id:[0-9]+}", ph.GetProduct)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getr.Handle("/docs", sh)
	getr.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	putr := r.Methods(http.MethodPut).Subrouter()
	putr.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)
	putr.Use(v.ProductValidationMiddleware)

	postr := r.Methods(http.MethodPost).Subrouter()
	postr.HandleFunc("/", ph.AddProduct)
	postr.Use(v.ProductValidationMiddleware)

	delr := r.Methods(http.MethodDelete).Subrouter()
	delr.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

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
