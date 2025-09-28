package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

// SaltedHash returns the hex sha256 sum of the concetenated salt and
// nhs number byte values.
func SaltedHash(salt []byte, nhsNumber int) string {
	no := []byte(strconv.Itoa(nhsNumber))
	bytesForHash := append(salt, no...)
	return fmt.Sprintf("%x", sha256.Sum256(bytesForHash))
}
