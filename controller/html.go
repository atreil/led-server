package main

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
)

type TemplateOptions struct {
	LedOpts []*LEDOption
}

type LEDOption struct {
	Id, Name, Type string
}

func Serve() (string, error) {
	tmplFile, err := filepath.Abs("index.go.txt")
	if err != nil {
		return "", fmt.Errorf("failed to resolve path to 'index.go.txt': %v", err)
	}

	tmpl, err := template.New("index.go.txt").ParseFiles(tmplFile)
	if err != nil {
		return "", fmt.Errorf("failed to parse template file (%s): %v", tmplFile, err)
	}

	opts := &TemplateOptions{
		LedOpts: []*LEDOption{
			{
				Id:   "led_min_frequency",
				Name: "Min frequency:",
				Type: "number",
			},
			{
				Id:   "led_max_frequency",
				Name: "Max frequency:",
				Type: "number",
			},
			{
				Id:   "led_n_fft_bins",
				Name: "FFT Bins:",
				Type: "number",
			},
		},
	}

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, opts); err != nil {
		return "", fmt.Errorf("failed to execute template file with options (%+v): %v", opts, err)
	}

	return buf.String(), nil
}
