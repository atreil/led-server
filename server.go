package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

func handleClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Not allowed"))
		return
	}

	mux.Lock()
	defer mux.Unlock()

	leds := make([]uint32, *numLeds)
	if err := dev.SetLedsSync(0, leds); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

var dev *ws2811.WS2811
var mux *sync.Mutex

var numLeds = flag.Int("leds", 300, "Number of leds")
var port = flag.String("port", ":80", "The port")

func main() {
	fmt.Println("Starting up...")
	var err error
	mux = &sync.Mutex{}

	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = 0
	opt.Channels[0].LedCount = *numLeds

	dev, err = ws2811.MakeWS2811(&opt)
	checkError(err)
	checkError(dev.Init())
	defer dev.Fini()

	c := make(chan os.Signal)
	w := make(chan int)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		w <- 1
	}()

	http.HandleFunc("/clear", handleClear)
	fmt.Println("Listening...")
	http.ListenAndServe(*port, nil)
	<-w
	fmt.Println("Goodbye!")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
