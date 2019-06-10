package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	h                 = flag.Bool("help", false, "this help")
	paramProtoOut     = flag.String("proto_out", "", "output protobuf define path(*.proto)")
	paramPbBinaryOut  = flag.String("pbbinary_out", "", "output protobuf binary data path(*.bytes)")
	paramProtoVersion = flag.Int("protover", 2, "output .proto file version, 2 or 3")
	paramDirInputMode = flag.Bool("dirmode", false, "input is dir(default is filename)")
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: excelconvert [-proto_out proto out path] [-pbbinary_out pbbinary out path] [-protover proto version(2,3)]
	excel filepath
Options:
`)
	flag.PrintDefaults()
}
