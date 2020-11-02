package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/flagsibh/mservices/product-api/server"
	"github.com/flagsibh/mservices/product-images/files"
	"github.com/flagsibh/mservices/product-images/handlers"
	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
)

const (
	// FileNamePattern defines the URI pattern for uploading files.
	FileNamePattern = "/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}"
)

var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "./imagestore", "Base path to save images")

func main() {
	env.Parse()

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-images",
			Level: hclog.LevelFromString(*logLevel),
		},
	)

	// create a logger for the server from the default logger
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// create the storage class, use local storage
	// max filesize 5MB
	stor, err := files.NewLocal(*basePath, 1024*1000*5)
	if err != nil {
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	// create a new serve mux and register the handlers
	r := mux.NewRouter()

	// create the handlers
	fh := handlers.NewFiles(stor, sl)

	ph := r.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc(FileNamePattern, fh.Upload)
	ph.HandleFunc("/", fh.UploadMultipart)

	// get files
	gh := r.Methods(http.MethodGet).Subrouter()
	gh.Handle(
		FileNamePattern,
		http.StripPrefix("/images/", http.FileServer(http.Dir(*basePath))),
	)

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
