package util

import (
	"crypto/rsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

// Signer interface
type Signer interface {
	GetPrivateKey() *rsa.PrivateKey
}

// Check error function
func Check(format string, args ...interface{}) {
	var hasError bool
	for _, v := range args {
		if err, ok := v.(error); ok && err != nil {
			hasError = true
			break
		}
	}
	if hasError {
		panic(fmt.Errorf(format, args...))
	}
}

// Uint64ToBytes func
func Uint64ToBytes(v uint64) []byte {
	bb := make([]byte, 8)
	binary.LittleEndian.PutUint64(bb, v)
	return bb
}

// UintToBytes func
func UintToBytes(v uint) []byte {
	bb := make([]byte, 4)
	binary.LittleEndian.PutUint32(bb, uint32(v))
	return bb
}

// JSONMarshal func
func JSONMarshal(w io.Writer, v interface{}) {
	if d, err := json.Marshal(v); err == nil {
		w.Write(d)
	} else {
		panic(err)
	}
}

// JSONUnmarshal func
func JSONUnmarshal(b []byte, v interface{}) interface{} {
	if err := json.Unmarshal(b, v); err != nil {
		panic(err)
	}
	return v
}

// InStringSlice func
func InStringSlice(v string, list ...string) bool {
	for _, vl := range list {
		if v == vl {
			return true
		}
	}
	return false
}
