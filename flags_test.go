package main

import (
	"fmt"
	"os"
	"testing"
)

func TestParseFlags(t *testing.T) {

	tests := []struct {
		args []string
		err  bool
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
	}
	for ii, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", ii), func(t *testing.T) {
			os.Args = tt.args
			_, err := ParseFlags()
			if tt.err && err == nil {
				t.Error("expected error")
			}
			if !tt.err && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
