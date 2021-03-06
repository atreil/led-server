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

type Request struct {
	Command string
}

func handleDaemonCommand(w http.ResponseWriter, r *http.Request) error {
	payload := &Request{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	if _, ok := allowed[payload.Command]; !ok {
		return fmt.Errorf("invalid command: %v", payload.Command)
	}

	cmd := exec.Command("systemctl", payload.Command, ServiceName)
	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("error running '%v': %v", cmd.Args, string(exitError.Stderr))
		}
		return fmt.Errorf("could not run command '%v': %v", cmd.Args, err)
	}

	if _, err := w.Write(output); err != nil {
		return fmt.Errorf("failed to write output but command succeeded: %v", err)
	}
	return nil
}
