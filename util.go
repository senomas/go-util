package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
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
func Check(err error, args ...interface{}) {
	if err != nil {
		if len(args) == 0 {
			panic(err)
		} else {
			format := args[0].(string)
			args[0] = err
			panic(fmt.Errorf(format, append(args[1:], err)...))
		}
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

// KeyRSA func
func KeyRSA() *rsa.PrivateKey {
	var privateKey *rsa.PrivateKey
	var err error

	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		panic(err)
	}

	return privateKey
}

// MarshalPublicKey func
func MarshalPublicKey(key *rsa.PublicKey) []byte {
	if data, err := x509.MarshalPKIXPublicKey(key); err != nil {
		panic(err)
	} else {
		return data
	}
}

// UnmarshalPublicKey func
func UnmarshalPublicKey(key []byte) *rsa.PublicKey {
	var kk interface{}
	var err error
	if kk, err = x509.ParsePKIXPublicKey(key); err != nil {
		panic(err)
	}
	return kk.(*rsa.PublicKey)
}

// MarshalPrivateKey func
func MarshalPrivateKey(key *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(key)
}

// UnmarshalPrivateKey func
func UnmarshalPrivateKey(key []byte) *rsa.PrivateKey {
	var kk *rsa.PrivateKey
	var err error
	if kk, err = x509.ParsePKCS1PrivateKey(key); err != nil {
		panic(err)
	}
	return kk
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

// PrivateKey func
func PrivateKey() *rsa.PrivateKey {
	var err error
	var privateKey *rsa.PrivateKey

	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		panic(err)
	}

	return privateKey
}

// EncryptRSA func
func EncryptRSA(key *rsa.PublicKey, data []byte) []byte {
	var err error
	klen := key.N.BitLen()/8 - 11
	if len(data) <= klen {
		var bb []byte
		if bb, err = rsa.EncryptPKCS1v15(rand.Reader, key, data); err != nil {
			panic(err)
		}
		return bb
	}
	var buf bytes.Buffer
	var bb []byte
	for i, w, r := 0, 0, len(data); r > 0; i, r = i+w, r-w {
		if r <= klen {
			if bb, err = rsa.EncryptPKCS1v15(rand.Reader, key, data[i:]); err != nil {
				panic(err)
			}
			buf.Write(bb)
			w = r
		} else {
			if bb, err = rsa.EncryptPKCS1v15(rand.Reader, key, data[i:i+klen]); err != nil {
				panic(err)
			}
			buf.Write(bb)
			w = klen
		}
	}
	return buf.Bytes()
}

// DecryptRSA func
func DecryptRSA(key *rsa.PrivateKey, data []byte) []byte {
	var err error
	klen := key.N.BitLen() / 8
	if len(data) <= klen {
		var bb []byte
		if bb, err = rsa.DecryptPKCS1v15(rand.Reader, key, data); err != nil {
			panic(err)
		}
		return bb
	}
	var buf bytes.Buffer
	var bb []byte
	for i, w, r := 0, 0, len(data); r > 0; i, r = i+w, r-w {
		if r <= klen {
			if bb, err = rsa.DecryptPKCS1v15(rand.Reader, key, data[i:]); err != nil {
				panic(err)
			}
			buf.Write(bb)
			w = r
		} else {
			if bb, err = rsa.DecryptPKCS1v15(rand.Reader, key, data[i:i+klen]); err != nil {
				panic(err)
			}
			buf.Write(bb)
			w = klen
		}
	}
	return buf.Bytes()
}

// KeyAES func
func KeyAES() []byte {
	var err error
	var key = make([]byte, 32)
	if _, err = rand.Read(key); err != nil {
		panic(fmt.Errorf("Error generate aes key\n\n%+v", err))
	}
	return key
}

// EncryptAES func
func EncryptAES(key, data []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext
}

// DecryptAES func
func DecryptAES(key, ciphertext []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Errorf("Error creating new block cipher\n%v\n", err))
	}

	if len(ciphertext) < aes.BlockSize {
		panic(fmt.Errorf("ciphertext too short"))
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext
}
