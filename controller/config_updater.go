package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Stream interface {
	io.ReadWriteSeeker
	Truncate(size int64) error
}

type Config struct {
	mux              sync.Mutex
	configPath       string
	configPathStream Stream
	d                Daemon
}

// the json config
type UpdateRequest struct {
	N_FFT_BINS, MIN_FREQUENCY, MAX_FREQUENCY, VISUALIZATION_TYPE *int
	SPECTRUM_BASE                                                []int
}

func NewConfig(configPath string, d Daemon) (*Config, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(absPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file '%s': %v", absPath, err)
	}

	return &Config{
		configPath:       absPath,
		configPathStream: f,
		d:                d,
	}, nil
}

func (c *Config) Update(reqRaw []byte) error {
	req := &UpdateRequest{}
	err := json.Unmarshal(reqRaw, req)

	if err != nil {
		log.Printf("failed to unmarshal payload %s: %v", string(reqRaw), err)
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	reqRaw, err = json.Marshal(req)
	if err != nil {
		log.Printf("failed to marshal payload %v: %v", req, err)
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	_, err = c.d.HandleDaemonCommand(Request{"stop"})
	if err != nil {
		return err
	}

	err = c.write(reqRaw)
	if err != nil {
		return err
	}

	_, err = c.d.HandleDaemonCommand(Request{"start"})
	return err
}

func (c *Config) write(data []byte) error {
	_, err := c.configPathStream.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to set file pointer at 0: %v", err)
	}

	// TODO: handle case where data was partially written
	n, err := c.configPathStream.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data at path '%s': %v", c.configPath, err)
	}

	return c.configPathStream.Truncate(int64(n))
}

func (c *Config) makeHandleUpdateRequest() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		req := make([]byte, r.ContentLength)
		_, err := r.Body.Read(req)
		if err != nil {
			return err
		}
		return c.Update(req)
	}
}
