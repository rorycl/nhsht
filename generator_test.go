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

func TestGenerateNumber(t *testing.T) {

	tests := []struct {
		nr        []NHSNumberRange
		records   int
		verbose   bool
		earlyDone bool
		count     int // number of records expected
	}{
		{
			nr:        []NHSNumberRange{NHSNumberRange{100_000_001, 100_000_100}},
			records:   10,
			verbose:   false,
			earlyDone: false,
			count:     10,
		},
		{
			nr:        []NHSNumberRange{NHSNumberRange{100_000_001, 100_000_100}},
			records:   10,
			verbose:   true,
			earlyDone: false,
			count:     10,
		},
		{
			nr:        []NHSNumberRange{NHSNumberRange{100_000_001, 100_000_100}},
			records:   10000,
			verbose:   true,
			earlyDone: false,
			count:     100,
		},
		{
			nr:        []NHSNumberRange{NHSNumberRange{100_000_001, 100_000_100}},
			records:   10,
			verbose:   true,
			earlyDone: true,
			count:     0,
		},
	}
	for ii, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", ii), func(t *testing.T) {
			NumberRanges = tt.nr // override package NumberRanges
			done := make(chan struct{})

			gn := generateNumber(tt.records, tt.verbose, done)
			i := 0
			if tt.earlyDone {
				done <- struct{}{}
			}
			for range gn {
				i++
			}
			if got, want := i, tt.count; got != want {
				t.Errorf("got %d want %d", got, want)
			}
		})
	}

}
