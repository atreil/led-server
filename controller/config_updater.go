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
	Sync() error
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

func (c *Config) Update(req *UpdateRequest) error {
	reqRaw, err := json.Marshal(req)
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
		return fmt.Errorf("failed to update config: %v", err)
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
	if err == io.ErrShortWrite {
		return fmt.Errorf("failed to update config because of: %v. please try again", err)
	} else if err != nil {
		return fmt.Errorf("failed to write data at path '%s': %v", c.configPath, err)
	}

	if err := c.configPathStream.Truncate(int64(n)); err != nil {
		return fmt.Errorf("failed to truncate file (%v): %v", c.configPath, err)
	}

	return c.configPathStream.Sync()
}

func (c *Config) makeHandleUpdateRequest() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		payload := &UpdateRequest{}
		if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
			return err
		}
		return c.Update(payload)
	}
}
