package main

import (
	"fmt"
	"io"
	"os"

	"github.com/parquet-go/parquet-go"
)

// NHSNoHash is a struct for encoding a pair of values representing an
// NHS Number and a salted SHA256 hash of the NHSNo for serialising into
// parquet format. For efficiency the NHSNo is stored as an int32 which
// precludes storing NHS numbers with leading zeros.
type NHSNoHash struct {
	NHSNo int32  `parquet:"nhsno,zstd"` // not a string
	Hash  string `parquet:"hash,zstd"`
}

// parquetWriter provides an NHSNoHash chan for writing records to a
// parquet file, providing also an error chan for reporting writing or
// other errors.
//
// Usage example:
//
//  // Initialiase parquetWriter.
//	writerChan, errChan, err := parquetWriter(fn)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// Write records to the parquet file.
//	for i, rec := range testRecords {
//		select {
//		case writerChan <- rec: // check writer is still open
//		case err := <-errChan: // fail if error
//			// handle error
//		}
//	}
//
//	// Signal completion of writing of records to parquetWriter.
//	close(writerChan)
//
//	// Synchronisation point to wait for flushing to disk.
//	if err := <-errChan; err != nil {
//		t.Fatal(err)
//	}

func parquetWriter(filename string) (chan<- NHSNoHash, <-chan error, error) {

	var err error
	writerChan := make(chan NHSNoHash)
	errorChan := make(chan error, 1)

	f, err := os.Create(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("parquet file creation error: %w", err)
	}

	writer := parquet.NewWriter(f)

	go func() {
		defer close(errorChan)
		i := 0
		for record := range writerChan {
			if err = writer.Write(record); err != nil {
				errorChan <- fmt.Errorf("parquet file row %d write error: %w", i, err)
				return
			}
			i++
		}
		// close parquet writer and flush file to disk
		if err = writer.Close(); err != nil {
			errorChan <- fmt.Errorf("parquet writer close error: %w", err)
			return
		}
		if err = f.Close(); err != nil {
			errorChan <- fmt.Errorf("file close error: %w", err)
			return
		}
	}()

	return writerChan, errorChan, err
}

// parquetReader reads NHSNoHash records from a parquet file. This is
// intended for testing purposes. To lookup an NHS number or hash using
// duckdb is recommended.
func parquetReader(filename string) ([]NHSNoHash, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	pq := parquet.NewReader(f)
	nhsHashes := []NHSNoHash{}
	i := 0
	for {
		var nh NHSNoHash
		err := pq.Read(&nh)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parquet row %d reading error: %w", i, err)
		}
		nhsHashes = append(nhsHashes, nh)
		i++
	}
	return nhsHashes, nil
}
