package main

import (
	"fmt"
	"testing"
)

// TestSaltedHash uses pre-generated values. Any change to a test
// 'expected' value should be carefully investigated.
func TestSaltedHash(t *testing.T) {
	tests := []struct {
		salt     []byte
		no       int
		expected string
	}{
		{
			salt:     []byte("456"),
			no:       123,
			expected: "c1cf024576e9c756b252bd5035efc64c72c17affe236909ded190d266a5bfdf1",
		},
		{
			salt:     []byte("5400fc717d9b2543f5e24da4b2c52f196845455073fd7fcef704c792322a552c"),
			no:       9000000000,
			expected: "3470703415c664edfb260d09336903a78650d28950441200cac09d95f09a9a7c",
		},
	}
	for ii, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", ii), func(t *testing.T) {
			if got, want := SaltedHash(tt.salt, tt.no), tt.expected; got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}
