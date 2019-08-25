package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type handlerFunc func(http.ResponseWriter, *http.Request) error

func makeHandler(fn handlerFunc, protocols ...string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		allowed := false
		for _, protocol := range protocols {
			if r.Method == protocol {
				allowed = true
				break
			}
		}

		if !allowed {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err := fn(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v\n", err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := dev.Render(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

var port = flag.String("port", ":8080", "The port")

func newHTTPServer() *http.Server {
	mux := http.NewServeMux()

	// Setup handlers
	mux.HandleFunc("/daemon", makeHandler(handleDaemonCommand, http.MethodPost))

	return &http.Server{
		Addr:           *port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func main() {
	log.Println("starting up...")

	ctx := context.Background()
	server := newHTTPServer()
	defer server.Shutdown(ctx)
	server.RegisterOnShutdown(func() {
		log.Printf("shutting down")
	})

	log.Println("listening...")
	err := server.ListenAndServe()
	log.Printf("shutting down with error: %v", err)
}
