package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Salt is representation of a 256bit randomly generated salt.
type Salt []byte

// NewSalt generates a new salt.
func NewSalt() (Salt, error) {
	saltBytes := make([]byte, 32)
	_, err := rand.Read(saltBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random salt: %w", err)
	}
	return Salt(saltBytes), nil
}

// AsHex provides a hex string representation of the salt.
func (s Salt) AsHex() string {
	return hex.EncodeToString(s)
}

// SaveToFile saves the salt hex string representation to a file.
func (s Salt) SaveToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("file create error: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()
	_, err = fmt.Fprint(f, s.AsHex())
	if err != nil {
		return fmt.Errorf("file write error: %w", err)
	}
	return nil
}

// SaltFromFile loads a hex encloded 256 bit salt from file.
func SaltFromFile(filename string) (Salt, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("file open error: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()
	saltBytes := make([]byte, 64)
	salt := make([]byte, 32)
	_, err = io.ReadAtLeast(f, saltBytes, 64)
	if err != nil {
		return nil, fmt.Errorf("salt reading error: %w", err)
	}
	_, err = hex.Decode(salt, saltBytes)
	if err != nil {
		return nil, fmt.Errorf("hex decoding error: %w", err)
	}
	return salt, nil
}
