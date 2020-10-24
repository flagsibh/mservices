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
)

func main() {
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)
	ph := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	srv := server.New(sm)
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
