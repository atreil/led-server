package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var configPath = flag.String("config_path", "", "Path to the config for the led visualizer")

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
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("%v\n", err.Error())))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func makeHandleDeviceCommand(dev *Device) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		payload := &Request{}
		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %v", err)
		}

		if payload.Command != strings.ToLower("clear") {
			return fmt.Errorf("expected command 'clear' but got %q", payload.Command)
		}

		return dev.Clear()
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) error {
	index, err := Serve()
	if err != nil {
		return err
	}
	w.Write([]byte(index))
	return nil
}

var port = flag.String("port", ":8080", "The port")

func newHTTPServer(dev *Device, daemon *DefaultDaemon, config *Config) (*http.Server, error) {
	mux := http.NewServeMux()

	// Setup handlers
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	log.Printf("setting file server root: %v", root)
	mux.HandleFunc("/", makeHandler(serveIndex, http.MethodGet))
	mux.HandleFunc("/daemon", makeHandler(daemon.MakeHandleDaemonCommandRequest(), http.MethodPost))
	mux.HandleFunc("/device", makeHandler(makeHandleDeviceCommand(dev), http.MethodPost))
	mux.HandleFunc("/led", makeHandler(config.makeHandleUpdateRequest(), http.MethodPost))

	return &http.Server{
		Addr:           *port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}, nil
}

func main() {
	log.Println("starting up...")

	flag.Parse()
	ctx := context.Background()

	if *configPath == "" {
		log.Panic("flag '--config_path' must be set")
	}

	daemon := &DefaultDaemon{}

	config, err := NewConfig(*configPath, daemon)
	if err != nil {
		log.Panicf("failed to start LED package: %v", err)
	}

	ws2811Dev, cleanup, err := NewDevice()
	if err != nil {
		log.Panicf("failed to intialize ws2811 device: %v", err)
	}
	defer cleanup()

	server, err := newHTTPServer(ws2811Dev, daemon, config)
	if err != nil {
		log.Panicf("failed to start up http server: %v", err)
	}
	defer server.Shutdown(ctx)
	server.RegisterOnShutdown(func() {
		log.Printf("shutting down")
		cleanup()
	})

	log.Println("listening...")
	err = server.ListenAndServe()
	log.Printf("shutting down with error: %v", err)
}
