package util

import (
	"os"
	"strconv"
)

// PartWriter interface
type PartWriter interface {
	Next()

	Write(p []byte) (n int, err error)

	Close() error

	Part() uint64
}

// PartWriterData struct
type PartWriterData struct {
	fn   string
	part uint64
	file *os.File
}

// NewPartWriter func
func NewPartWriter(path string, id uint64) *PartWriterData {
	var err error
	var file *os.File
	fn := path + "/F" + strconv.FormatUint(id, 16) + "."
	if file, err = os.Create(fn + strconv.FormatUint(uint64(1), 16) + ".tmp"); err != nil {
		panic(err)
	}
	return &PartWriterData{fn, 1, file}
}

// Next func
func (w *PartWriterData) Next() {
	var err error
	if err = w.file.Close(); err != nil {
		panic(err)
	}
	w.part++
	if w.file, err = os.Create(w.fn + strconv.FormatUint(uint64(w.part), 16) + ".tmp"); err != nil {
		panic(err)
	}
}

// Write func
func (w *PartWriterData) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}

// Close func
func (w *PartWriterData) Close() error {
	return w.file.Close()
}

// Part func
func (w *PartWriterData) Part() uint64 {
	return w.part
}
