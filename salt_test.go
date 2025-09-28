package main

import (
	"os"
	"testing"
)

func TestSalt(t *testing.T) {

	// Setup a test output parquet file.
	tmpFile, err := os.CreateTemp("", "tpw_*.salt")
	if err != nil {
		t.Fatal(err)
	}
	fn := tmpFile.Name()

	t.Cleanup(func() {
		_ = os.Remove(fn) // remove the test file
	})

	salt, err := NewSalt()
	if err != nil {
		t.Fatal(err)
	}

	err = salt.SaveToFile(fn)
	if err != nil {
		t.Fatal(err)
	}

	// Read salt written to file matches the original salt.
	fileSalt, err := SaltFromFile(fn)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := fileSalt.AsHex(), salt.AsHex(); got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
