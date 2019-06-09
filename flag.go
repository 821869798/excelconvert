package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	h                 = flag.Bool("help", false, "this help")
	paramPackageName  = flag.String("package", "", "set the package name in table")
	paramTableName    = flag.String("tname", "table", "set the table name")
	paramProtoOut     = flag.String("proto_out", "", "output protobuf define (*.proto)")
	paramPbBinaryOut  = flag.String("pbbinary_out", "", "output protobuf binary data (*.bytes)")
	paramProtoVersion = flag.Int("protover", 2, "output .proto file version, 2 or 3")
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: excelconvert [-package packagename] [-tname tablename] [-proto_out proto out path] [-pbbinary_out pbbinary out path] [-protover proto version(2,3)]
	excel filepath
Options:
`)
	flag.PrintDefaults()
}
