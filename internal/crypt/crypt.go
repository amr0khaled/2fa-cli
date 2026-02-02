package crypt

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash"
)

const SIZE int = 32

func randomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func R(data []byte, key []byte) []byte {
	cipher := hmac.New(func() hash.Hash {
		sha := sha256.New()
		sha.Write(data)
		return sha
	}, key)
	return cipher.Sum(nil)
}

func Gen() []byte {
	random, err := randomBytes(32)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	return random
}
