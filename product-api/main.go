package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/flagsibh/mservices/product-api/handlers"
	mw "github.com/flagsibh/mservices/product-api/handlers/middleware"
	"github.com/flagsibh/mservices/server"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
)

const (
	// IDParamPattern sets the Regex pattern to extract de ID from de URI.
	IDParamPattern = "/{id:[0-9]+}"
)

func main() {
	l := hclog.Default().Named("product-api")
	v := mw.NewValidation(l)

	ph := handlers.NewProducts(l)

	r := mux.NewRouter()
	r.Use(mw.ContentTypeMiddleware)

	getr := r.Methods(http.MethodGet).Subrouter()
	getr.HandleFunc("/", ph.GetProducts)
	getr.HandleFunc(IDParamPattern, ph.GetProduct)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)
	getr.Handle("/docs", sh)
	getr.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	putr := r.Methods(http.MethodPut).Subrouter()
	putr.HandleFunc(IDParamPattern, ph.UpdateProduct)
	putr.Use(v.ProductValidationMiddleware)

	postr := r.Methods(http.MethodPost).Subrouter()
	postr.HandleFunc("/", ph.AddProduct)
	postr.Use(v.ProductValidationMiddleware)

	delr := r.Methods(http.MethodDelete).Subrouter()
	delr.HandleFunc(IDParamPattern, ph.DeleteProduct)

	srv := server.New(r, l)
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			l.Info("Shutdown", "msg", err)
		}
	}()

	channel := make(chan os.Signal)
	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	sc := <-channel
	l.Info("Received terminate, graceful shutdown", "channel", sc)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(tc)
}
