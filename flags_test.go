package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {

	tests := []struct {
		args    []string
		verbose bool
		err     bool
	}{
		{
			args: []string{"program", "-s", "salt", "-p", "file.parquet"},
			err:  false,
		},
		{
			args: []string{"program", "-s", "salt", "-p", "file.parquet", "-r", "20"},
			err:  false,
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
			args:    []string{"program", "-s", "salt", "-p", "file.parquet", "-v"},
			verbose: true,
			err:     false,
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
				return
			}
			if tt.err == false && err != nil {
				t.Fatalf("unexpected error: %v", err)
				return
			}
			if got, want := f.Verbose, tt.verbose; got != want {
				t.Errorf("verbose got %t want %t", got, want)
			}
		})
	}
}
