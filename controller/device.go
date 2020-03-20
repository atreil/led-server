package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

var globalDev *Device

type Device struct {
	// The device interface
	dev *ws2811.WS2811
	mux sync.Mutex
}

func (d *Device) Clear() error {
	d.mux.Lock()
	defer d.mux.Unlock()

	leds := make([]uint32, *numLeds)
	if err := d.dev.SetLedsSync(0, leds); err != nil {
		return fmt.Errorf("error setting leds: %s", err)
	}

	return nil
}

func NewDevice() (*Device, func() error, error) {
	if globalDev != nil {
		return globalDev, globalDev.cleanup, nil
	}

	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = 0
	opt.Channels[0].LedCount = *numLeds

	dev, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		log.Printf("failed to create WS2811 device: %v", err)
		return nil, nil, err
	}

	if err = dev.Init(); err != nil {
		log.Printf("failed to initalize WS2811 device: %v", err)
		return nil, nil, err
	}

	globalDev := &Device{dev: dev}
	return globalDev, globalDev.cleanup, nil
}

func (d *Device) cleanup() error {
	if d == nil || globalDev == nil {
		return nil
	}

	d.mux.Lock()
	defer d.mux.Unlock()
	log.Println("cleaning up device")
	d.dev.Fini()
	globalDev = nil
	return nil
}

var numLeds = flag.Int("leds", 300, "Number of leds")
