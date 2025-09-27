package main

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestParquetWriter(t *testing.T) {

	tmpFile, err := os.CreateTemp("", "tpw.parquet")
	if err != nil {
		t.Fatal(err)
	}
	fn := tmpFile.Name()
	// defer func() {
	// 	_ = os.Remove(fn)
	// }()

	writeChan := make(chan NHSNoHash)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := parquetWriter(fn, writeChan)
		if err != nil {
			t.Fatal(err)
		}
		wg.Done()
	}()

	for i, rec := range []NHSNoHash{
		NHSNoHash{300_000_000, "98ea6e4f216f2fb4b69fff9b3a44842c38686ca685f3f55dc48c5d3fb1107be4"},
		NHSNoHash{300_000_001, "c641344867e9806fadfd219f25b62b97c94db0eed04a1d79e93676533cfb782b"},
	} {
		fmt.Printf("writing record %d: %v\n", i, rec)
		writeChan <- rec
	}

	close(writeChan)
	wg.Wait()

}
