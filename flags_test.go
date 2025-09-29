package main

import (
	"fmt"
	"os"
	"runtime"
	"testing"
)

func TestParseFlags(t *testing.T) {

	defaultNoGoroutines := uint(runtime.NumCPU() * GoRoutineNumCPUFactor)

	tests := []struct {
		args        []string
		verbose     bool
		numRoutines uint
		err         bool
	}{
		{
			args:        []string{"program", "-s", "salt", "-p", "file.parquet"},
			err:         false,
			numRoutines: defaultNoGoroutines,
		},
		{
			args:        []string{"program", "-s", "salt", "-p", "file.parquet", "-r", "20"},
			err:         false,
			numRoutines: defaultNoGoroutines,
		},
		{
			args: []string{"-p", "file.parquet", "-r", "20"},
			err:  true,
		},
		{
			args: []string{"program", "-s", "salt", "-r", "20"},
			err:  true,
		},
		{
			args:        []string{"program", "-s", "salt", "-p", "file.parquet", "-v"},
			verbose:     true,
			err:         false,
			numRoutines: defaultNoGoroutines,
		},
		{
			args: []string{"program", "-s", "salt", "-p", "file.parquet", "-g", "'-3'"},
			err:  true,
		},
		{
			args:        []string{"program", "-s", "salt", "-p", "file.parquet", "-g", "40"},
			err:         false,
			numRoutines: 40,
		},
		{
			args: []string{"program", "-s", "salt", "-p", "file.parquet", "-g", "800000"},
			err:  true,
		},
	}
	for ii, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", ii), func(t *testing.T) {
			os.Args = tt.args
			f, err := ParseFlags()
			if tt.err {
				if err == nil {
					t.Error("expected error")
					return
				}
				fmt.Println(err)
				return
			}
			if tt.err == false && err != nil {
				fmt.Println(err)
				t.Fatalf("unexpected error: %v", err)
				return
			}
			if got, want := f.Verbose, tt.verbose; got != want {
				t.Errorf("verbose got %t want %t", got, want)
			}
			if got, want := f.GoRoutines, tt.numRoutines; got != want {
				t.Errorf("number of goroutines got %d want %d", got, want)
			}
		})
	}
}
