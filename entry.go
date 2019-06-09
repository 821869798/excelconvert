package main

import (
	"flag"
	"fmt"
	"github.com/821869798/excelconvert/converter"
	"github.com/821869798/excelconvert/excel"
	_ "github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	_ "github.com/jhump/protoreflect/desc/builder"
	_ "github.com/jhump/protoreflect/desc/protoparse"
	_ "github.com/tealeg/xlsx"
	"os"
)

func Entry() {
	g := converter.NewGlobals()

	g.ProtoVersion = *paramProtoVersion

	g.InputFileList = flag.Arg(0)

	if *paramProtoOut != "" && *paramPbBinaryOut != "" {
		g.AddOutputType("proto", &converter.ProtoArgs{
			PackageName: *paramPackageName,
			TableName:   *paramTableName,
			ProtoOut:    *paramProtoOut,
			PbBinaryOut: *paramPbBinaryOut,
		})
	}

	if len(g.Converters) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "\nnot output file set!")
		return
	}

	if !StartExport(g) {
		os.Exit(1)
	}
}

func StartExport(g *converter.Globals) bool {

	glog.Info("==========%s==========", "开始导出")
	//g.InputFileList

	file := excel.NewFile(g.InputFileList)
	if file == nil {
		return false
	}

	file.GlobalFD = g.FileDescriptor

	// 电子表格数据导出到Table对象
	if !file.ExportLocalType() {
		return false
	}

	return false
}
