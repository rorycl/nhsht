package main

import (
	"fmt"
	"runtime"

	"github.com/jessevdk/go-flags"
)

const version string = "0.0.1"

// Flags are the flag options consuming os.Args.
type Flags struct {
	SaltFile    string `short:"s" long:"saltfile" description:"file to save hex encoded salt (required)" required:"yes"`
	ParquetFile string `short:"p" long:"parquetfile" description:"hash table parquet file path (required)" required:"yes"`
	Records     uint   `short:"r" long:"records" description:"only generate this number of records"`
	GoRoutines  uint   `short:"g" long:"goroutines" description:"number of goroutines (default 8 * numcpu)"`
	Verbose     bool   `short:"v" long:"verbose" description:"report progress of number generation"`
}

// Number of goroutines to run by default as a factor of the CPU
const GoRoutineNumCPUFactor int = 8

// Maximum number of goroutines to run
const GoRoutineNumCPUFactorMax int = 30

var cmdTpl string = `

NHS Number salted hash table generator.

version %s

This program: 

* generates a random 256bit salt and saves this, hex encoded, to the
  specified file.

* generates all of the currently possible NHS numbers for England and
  saves each number with a salted sha256 value in the specified parquet
  file which is zstd compressed. Not that a full hash table parquet file
  is likely to be around 20GB in size.

Use the "records" flag to avoid generating all possible 300 million NHS
numbers for England.

e.g.
  nhsht -s salt -p hashes.parquet
to generate only 20 records:
  nhsht -s salt -p hashes.parquet -r 20
to show progress:
  nhsht -s salt -p hashes.parquet -v`

// ParseFlags parses the command line options.
func ParseFlags() (*Flags, error) {
	var options Flags
	var parser = flags.NewParser(&options, flags.HelpFlag)
	parser.Usage = fmt.Sprintf(cmdTpl, version)

	if _, err := parser.Parse(); err != nil {
		return nil, err
	}

	if options.GoRoutines == 0 {
		options.GoRoutines = uint(runtime.NumCPU() * GoRoutineNumCPUFactor)
	}
	if options.GoRoutines > uint(runtime.NumCPU()*GoRoutineNumCPUFactorMax) {
		err := fmt.Errorf(
			"goroutine number of %d more than permitted maximum of %d",
			options.GoRoutines,
			GoRoutineNumCPUFactorMax,
		)
		return nil, err
	}

	return &options, nil
}
