package main

import (
	"fmt"
	"os"
)

func main() {

	flags, err := ParseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = Generator(flags.SaltFile, flags.ParquetFile, int(flags.Records), int(flags.GoRoutines), flags.Verbose)
	if err != nil {
		fmt.Println("generator error:", err)
		os.Exit(1)
	}

}
