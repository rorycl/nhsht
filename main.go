package main

import (
	"fmt"
	"os"
)

func main() {

	flags, err := ParseFlags()
	if err != nil {
		// jessevdk/go-flags emits its own error message
		os.Exit(1)
	}

	err = Generator(flags.SaltFile, flags.ParquetFile, int(flags.Records), flags.Verbose)
	if err != nil {
		fmt.Println("generator error:", err)
		os.Exit(1)
	}

}
