package main

import (
	"flag"
)

func main() {
	flag.Parse()

	if *h {
		usage()
		return
	}

	Entry()
}
