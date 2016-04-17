package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

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
