package main

import (
	"fmt"

	"github.com/jessevdk/go-flags"
)

const version string = "0.0.1"

// Flags are the flag options consuming os.Args.
type Flags struct {
	SaltFile    string `short:"s" long:"saltfile" description:"file to save hex encoded salt" required:"yes"`
	ParquetFile string `short:"p" long:"parquetfile" description:"hash table parquet file path" required:"yes"`
	Records     uint   `short:"r" long:"records" description:"only generate this number of records"`
}

var cmdTpl string = `

NHS Number salted hash table generator.

version %s

This program: 

* generates a random 256bit salt and saves this, hex
  encoded, to the specified file.

* generates all of the currently possible NHS numbers for England and
  saves each number with a salted sha256 value in a parquet file.

If you do not want to generate all of the possible 300 million NHS
numbers for England, perhaps for testing, provide the "records" flag
with an appropriate number.

e.g.
  nhsht -s salt -p hashes.parquet
or
  nhsht -s salt -p hashes.parquet -r 20

`

// ParseFlags parses the command line options.
func ParseFlags() (*Flags, error) {
	var options Flags
	var parser = flags.NewParser(&options, flags.Default)
	parser.Usage = fmt.Sprintf(cmdTpl, version)

	if _, err := parser.Parse(); err != nil {
		return nil, err
	}
	return &options, nil
}
