package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestModulo11(t *testing.T) {

	cases := []struct {
		name   string
		input  int
		output int
		isErr  bool
		err    error
	}{
		{
			name:   "too short",
			input:  94347659,
			output: -1,
			isErr:  true,
		},
		{
			name:   "too long",
			input:  9434765911,
			output: -1,
			isErr:  true,
		},
		{
			name:   "ok (remainder 11)",
			input:  943476590,
			output: 9434765900,
			isErr:  false,
		},
		{
			name:   "ok",
			input:  943476591,
			output: 9434765919,
			isErr:  false,
		},
		{
			name:   "ok",
			input:  943476592,
			output: 9434765927,
			isErr:  false,
		},
		{
			name:   "ok",
			input:  943476593,
			output: 9434765935,
			isErr:  false,
		},
		{
			name:   "ok",
			input:  943476594,
			output: 9434765943,
			isErr:  false,
		},
		{
			name:   "ok",
			input:  943476595,
			output: 9434765951,
			isErr:  false,
		},
		{
			name:   "invalid nhs number (checksum 10)",
			input:  943476596,
			output: -1,
			isErr:  true,
			err:    ErrInvalidModulo11,
		},
		{
			name:   "ok",
			input:  943476597,
			output: 9434765978,
			isErr:  false,
		},
		{
			name:   "ok",
			input:  943476598,
			output: 9434765986,
			isErr:  false,
		},
		{
			name:   "ok",
			input:  943476599,
			output: 9434765994,
			isErr:  false,
		},
	}

	for ii, tt := range cases {
		t.Run(fmt.Sprintf("%d_%s", ii, tt.name), func(t *testing.T) {
			output, err := modulo11(tt.input)
			if tt.isErr {
				if err == nil {
					t.Fatal("expected error")
				}
				if tt.err != nil {
					if !errors.Is(err, tt.err) {
						t.Fatalf("expected error '%T' got '%v'", tt.err, err)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			fmt.Printf("%d -> %d\n", tt.input, output)
			if got, want := output, tt.output; got != want {
				t.Errorf("got %d want %d", got, want)
			}
		})
	}
}

// TestModuloSequence tests that any sequence of 11 consecutive 9 digit
// numbers in the range has exactly one invalid ("checksum 10") value.
func TestModuloSequence(t *testing.T) {
	// start := 943476597
	start := 900000000
	invalidModuloCounter := 0
	loops := 0
	for i := 0; i <= 10; i++ {
		loops++
		_, err := modulo11(start + i)
		if err != nil && errors.Is(err, ErrInvalidModulo11) {
			invalidModuloCounter++
		}
	}
	if got, want := loops, 11; got != want {
		t.Errorf("loops got %d want %d", got, want)
	}
	if got, want := invalidModuloCounter, 1; got != want {
		t.Errorf("invalid counter got %d want %d", got, want)
	}
}

/*
func BenchmarkModuloFuncs(b *testing.B) {
	for _, tt := range []struct {
		name    string
		modFunc func(int) (int, error)
	}{
		{"modulo11", modulo11},
		{"modulo11b", modulo11b},
	} {
		b.Run(tt.name, func(b *testing.B) {
			input := 900000000
			for i := 0; i < b.N; i++ {
				_, _ = tt.modFunc(input)
				input++
			}
		})
	}
}
*/
