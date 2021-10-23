package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type FakeStream struct {
	idx  int64
	data []byte
}

func (f *FakeStream) Read(p []byte) (int, error) {
	return 0, errors.New("Read: not implemented")
}

func (f *FakeStream) Write(p []byte) (int, error) {
	pCopy := make([]byte, len(p))
	written := copy(pCopy, p)
	f.data = append(f.data[:f.idx], pCopy...)
	f.idx += int64(written)
	return written, nil
}

func (f *FakeStream) Seek(offset int64, whence int) (int64, error) {
	if whence != 0 {
		return 0, fmt.Errorf("whence (%v) must be 0", whence)
	}
	f.idx = offset
	return offset, nil
}

func (f *FakeStream) Truncate(size int64) error {
	if size < int64(cap(f.data)) {
		f.data = f.data[:size]
		return nil
	}

	toFill := size - int64(len(f.data))
	filler := make([]byte, toFill)
	f.data = append(f.data, filler...)
	return nil
}

func (f *FakeStream) Sync() error {
	return nil
}

type FakeDaemon struct {
	requestQ []Request
}

func (f *FakeDaemon) HandleDaemonCommand(req Request) ([]byte, error) {
	f.requestQ = append(f.requestQ, req)
	return nil, nil
}

func ptrInt(i int) *int {
	return &i
}

func TestConfig(t *testing.T) {
	dataInit := UpdateRequest{
		N_FFT_BINS: ptrInt(0),
	}
	dataInitRaw, err := json.Marshal(dataInit)
	if err != nil {
		t.Fatalf("setup failed marshalling data: %v", err)
	}

	fs := &FakeStream{
		data: dataInitRaw,
	}
	fd := &FakeDaemon{}
	c := &Config{
		configPath:       "test",
		configPathStream: fs,
		d:                fd,
	}

	wantData := &UpdateRequest{
		N_FFT_BINS: ptrInt(1),
	}
	wantDataRaw, err := json.Marshal(wantData)
	if err != nil {
		t.Fatalf("setup failed marshalling data: %v", err)
	}

	if err := c.Update(wantData); err != nil {
		t.Errorf("Update(%s) returned an unexpected error: %v", wantDataRaw, err)
	}

	if !bytes.Equal(fs.data, wantDataRaw) {
		t.Errorf("found unexpected data written to stream (got: %s, want: %s)", fs.data, wantDataRaw)
	}

	wantRequestQ := []Request{{"stop"}, {"start"}}
	if !reflect.DeepEqual(fd.requestQ, wantRequestQ) {
		t.Errorf("invalid order of requests (got: %v, want: %v)", fd.requestQ, wantRequestQ)
	}
}
