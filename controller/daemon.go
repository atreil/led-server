package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

const (
	ServiceName = "audio-reactive-led-strip.service"
)

var allowed = map[string]bool{
	"start":   true,
	"stop":    true,
	"restart": true,
	"enable":  true,
	"disable": true,
	"status":  true,
}

type Daemon interface {
	HandleDaemonCommand(req Request) ([]byte, error)
}

type DefaultDaemon struct {
}

type Request struct {
	Command string
}

func (d *DefaultDaemon) MakeHandleDaemonCommandRequest() func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		payload := Request{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			return fmt.Errorf("failed to unmarshal payload: %v", err)
		}
		output, err := d.HandleDaemonCommand(payload)
		if err != nil {
			return err
		}
		if _, err := w.Write(output); err != nil {
			return fmt.Errorf("failed to write output but command succeeded: %v", err)
		}
		return nil
	}
}

func (d *DefaultDaemon) HandleDaemonCommand(req Request) ([]byte, error) {
	if _, ok := allowed[req.Command]; !ok {
		return nil, fmt.Errorf("invalid command: %v", req.Command)
	}

	cmd := exec.Command("systemctl", req.Command, ServiceName)
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("error running '%v': %v", cmd.Args, string(exitError.Stderr))
		}
		return nil, fmt.Errorf("could not run command '%v': %v", cmd.Args, err)
	}
	return output, nil
}
