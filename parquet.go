package main

import (
	"fmt"
	"log"
	"os"

	"github.com/parquet-go/parquet-go"
)

// NHSNoHash is a struct for encoding a pair of values representing an
// NHS Number and a salted and SHA256 hash of the NHSNo for serialising
// into parquet format.
type NHSNoHash struct {
	NHSNo int32  `parquet:"nhsno,zstd"`
	Hash  string `parquet:"hash,zstd"`
}

// parquetWriter writes NHSNoHash records to a parquet file by reading
// from the readChan channel of NHSNoHash values.
func parquetWriter(filename string, readChan <-chan NHSNoHash) error {
	var err error

	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("parquet file creation error: %w", err)
	}

	writer := parquet.NewWriter(f)

	defer func() {
		if err = writer.Close(); err != nil {
			log.Printf("parquet file close error: %v", err)
		}
		if err = f.Close(); err != nil {
			log.Printf("file close error: %v", err)
		}
		fmt.Println("done!")
	}()

	for record := range readChan {
		if err = writer.Write(record); err != nil {
			return err
		}
		fmt.Printf("wrote %#v\n", record)
	}

	return err
}
