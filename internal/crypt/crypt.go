package crypt

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
)

type HashMethod interface {
	Hmac(data []byte, key []byte) []byte
	GenSecret() ([]byte, error)
}

type BaseHash struct {
	Size int
}

type Sha1 struct {
	BaseHash
}
type Sha256 struct {
	BaseHash
}
type Sha512 struct {
	BaseHash
}

func (base BaseHash) GenSecret() ([]byte, error) {
	bytes := make([]byte, base.Size)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func NewSha1() *Sha1 {
	return &Sha1{
		BaseHash: BaseHash{
			Size: 20,
		},
	}
}
func NewSha256() *Sha256 {
	return &Sha256{
		BaseHash: BaseHash{
			Size: 32,
		},
	}
}
func NewSha512() *Sha512 {
	return &Sha512{
		BaseHash: BaseHash{
			Size: 64,
		},
	}
}

func (Sha1) Hmac(data []byte, key []byte) []byte {
	cipher := hmac.New(sha1.New, key)
	cipher.Write(data)
	return cipher.Sum(nil)
}
func (Sha256) Hmac(data []byte, key []byte) []byte {
	cipher := hmac.New(sha256.New, key)
	cipher.Write(data)
	return cipher.Sum(nil)
}
func (Sha512) Hmac(data []byte, key []byte) []byte {
	cipher := hmac.New(sha512.New, key)
	cipher.Write(data)
	return cipher.Sum(nil)
}

// TODO: For CLI adaption
var Hashes = map[string]HashMethod{
	"sha1":   NewSha1(),
	"sha256": NewSha256(),
	"sha512": NewSha512(),
}

func RandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
