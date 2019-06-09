package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()

	if *h {
		usage()
		return
	}

	if flag.NArg() == 0 {
		usage()
		_, _ = fmt.Fprintln(os.Stderr, "\nnot a excel file input!")
		return
	}

	Entry()
}
