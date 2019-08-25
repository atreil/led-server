package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

type device struct {
	// The device interface
	dev *ws2811.WS2811
	mux sync.Mutex
}

func (d *device) handleClear(w http.ResponseWriter, r *http.Request) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	leds := make([]uint32, *numLeds)
	if err := d.dev.SetLedsSync(0, leds); err != nil {
		return fmt.Errorf("error setting leds: %s", err)
	}

	return nil
}

func newDevice() (*device, error) {
	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = 0
	opt.Channels[0].LedCount = *numLeds

	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		log.Printf("failed to create WS2811 device: v", err)
		return nil, err
	}

	if err = dev.Init(); err != nil {
		log.Printf("failed to initalize WS2811 device: %v", err)
		return nil, err
	}

	return &device{
		dev: dev,
	}, nil
}

func (d *device) Cleanup() {
	if d == nil {
		return
	}

	d.mux.Lock()
	defer d.mux.Unlock()
	log.Println("cleaning up device")
	d.dev.Fini()
}

var numLeds = flag.Int("leds", 300, "Number of leds")
var dev *device

// func main() {
// 	log.Println("starting up...")

// 	ctx := context.Background()
// 	var err error
// 	dev, err = newDevice()
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	server := newHTTPServer()
// 	defer server.Shutdown(ctx)
// 	server.RegisterOnShutdown(func() {
// 		dev.Cleanup()
// 	})

// 	log.Println("listening...")
// 	err = http.ListenAndServe(*port, nil)
// 	log.Printf("shutting down with error: %v", err)
// }
