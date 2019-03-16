package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

func handleClear(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Got method: %s\n", r.Method)
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Not allowed"))
		return
	}

	mux.Lock()
	defer mux.Unlock()

	var color int64
	if colorQuery, ok := r.URL.Query()["color"]; ok {
		var err error
		color, err = strconv.ParseInt(colorQuery[0], 0, 64)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		/*
			if color < 0 || color > 255 {
				fmt.Printf("Color must be between 0 and 255. Got %v\n", color)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		*/
	}

	leds := make([]uint32, *numLeds)
	if color != 0 {
		for i := 0; i < *numLeds; i++ {
			leds[i] = uint32(color)
		}
	}
	fmt.Printf("Using color: %v\n", color)
	fmt.Printf("Current leds: %v\n", dev.Leds(0))
	if err := dev.SetLedsSync(0, leds); err != nil {
		fmt.Printf("Error setting leds: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		fmt.Printf("Success\n")
		w.WriteHeader(http.StatusOK)
		if err := dev.Render(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

var dev *ws2811.WS2811
var mux *sync.Mutex

var numLeds = flag.Int("leds", 300, "Number of leds")
var port = flag.String("port", ":80", "The port")
var brightness = flag.Int("brightness", 100, "0-254")

func main() {
	fmt.Println("Starting up...")
	flag.Parse()
	var err error
	mux = &sync.Mutex{}

	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = *brightness
	opt.Channels[0].LedCount = *numLeds

	dev, err = ws2811.MakeWS2811(&opt)
	checkError(err)
	err = dev.Init()
	checkError(err)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("received signal")
		dev.Fini()
		fmt.Println("Buh bye")
		os.Exit(0)
	}()

	http.HandleFunc("/clear", handleClear)
	http.HandleFunc("/killkillkill", func(w http.ResponseWriter, r *http.Request) {
		c <- syscall.SIGTERM
	})
	fmt.Println("Listening...")
	http.ListenAndServe(*port, nil)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
