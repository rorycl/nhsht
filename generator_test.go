package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestGenerator(t *testing.T) {

	numberOfRecords := 200_000
	goroutineNo := runtime.NumCPU() * 8
	var err error

	// Setup a test output salt and parquet file.
	saltFile, err := os.CreateTemp("", "gensalt_*.sha256")
	if err != nil {
		t.Fatal(err)
	}
	sf := saltFile.Name()

	parquetFile, err := os.CreateTemp("", "genparquet_*.parquet")
	if err != nil {
		t.Fatal(err)
	}
	pqf := parquetFile.Name()

	t.Cleanup(func() {
		if err != nil {
			return // keep tmp files
		}
		_ = os.Remove(sf)
		_ = os.Remove(pqf)
	})

	start := time.Now()
	inMemory := false
	err = Generator(sf, pqf, numberOfRecords, goroutineNo, inMemory, true)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("generation time for %d records: %v\n",
		numberOfRecords, time.Since(start),
	)

	// Check number of records in parquet file.
	nhs, err := parquetReader(pqf)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(nhs), numberOfRecords; got != want {
		err = errors.New("invalid number of records")
		t.Errorf("got %d records want %d", got, want)
	}

}
