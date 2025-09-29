package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var testRecords []NHSNoHash = []NHSNoHash{
	NHSNoHash{300_000_000, "98ea6e4f216f2fb4b69fff9b3a44842c38686ca685f3f55dc48c5d3fb1107be4"},
	NHSNoHash{300_000_001, "c641344867e9806fadfd219f25b62b97c94db0eed04a1d79e93676533cfb782b"},
}

func TestParquetWriter(t *testing.T) {

	// Setup a test output parquet file.
	tmpFile, err := os.CreateTemp("", "tpw_*.parquet")
	if err != nil {
		t.Fatal(err)
	}
	fn := tmpFile.Name()

	t.Cleanup(func() {
		_ = os.Remove(fn) // remove the test file
	})

	inMemory := false
	writerChan, errChan, err := parquetWriter(fn, inMemory)
	if err != nil {
		t.Fatal(err)
	}

	// Write records to the parquet file.
	for i, rec := range testRecords {
		fmt.Printf("writing record %d: %v\n", i, rec)
		select {
		case writerChan <- rec: // check writer is still open
		case err := <-errChan: // fail if error
			t.Fatal(err)
		}
	}

	// Signal completion of writing of records to parquetWriter.
	close(writerChan)

	// Synchronisation point to wait for parquet file to be closed and
	// flushed to disk.
	if err := <-errChan; err != nil {
		t.Fatal(err)
	}

	// Check something was written to disk.
	stat, err := os.Stat(fn)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := stat.Size(), int64(500); got < want {
		t.Errorf("got file size %d want at least %d", got, want)
	}

	// Check contents of parquet file.
	nhs, err := parquetReader(fn)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(nhs), len(testRecords); got != want {
		t.Errorf("got %d records want %d", got, want)
	}

	if !cmp.Equal(nhs, testRecords) {
		t.Error(cmp.Diff(nhs, testRecords))
	}

}
