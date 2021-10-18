package main

import "sync"

type LED struct {
	mux sync.Mutex
}

type UpdateRequest struct {
	N_FFT_BINS, MIN_FREQUENCY, MAX_FREQUENCY, SPECTRUM_BASE, VISUALIZATION_TYPE string
}

type LocalConfig struct {
	N_FFT_BINS, MIN_FREQUENCY, MAX_FREQUENCY, SPECTRUM_BASE *int
	SpectrumBase                                            []int
}

func reqToConf(req *UpdateRequest) (*LocalConfig, error) {

}

func (l *LED) Update(req *UpdateRequest) error {

}
