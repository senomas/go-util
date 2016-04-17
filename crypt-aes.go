package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

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
