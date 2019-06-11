package main

import (
	"flag"
	"fmt"
	"github.com/821869798/excelconvert/converter"
	"github.com/821869798/excelconvert/excel"
	"github.com/821869798/excelconvert/model"
	"github.com/821869798/excelconvert/util"
	_ "github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	_ "github.com/jhump/protoreflect/desc/builder"
	_ "github.com/jhump/protoreflect/desc/protoparse"
	_ "github.com/tealeg/xlsx"
	"os"
	"path/filepath"
)

func Entry() {
	g := converter.NewGlobals()

	g.ProtoVersion = *paramProtoVersion

	//添加输入文件或者文件夹
	if *paramDirInputMode {
		for _, v := range flag.Args() {
			fileLists, err := util.GetExcellist(v)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "\n", v, ":not a correct dir input!")
				return
			}
			g.InputFileList = append(g.InputFileList, fileLists...)
		}
	} else {
		for _, v := range flag.Args() {
			g.InputFileList = append(g.InputFileList, v)
		}
	}

	//test
	//*paramProtoOut = `./example/protos`
	//*paramPbBinaryOut = `./example/bytes`
	//g.InputFileList = append(g.InputFileList, `./example/example.xlsx`)

	if len(g.InputFileList) == 0 {
		usage()
		_, _ = fmt.Fprintln(os.Stderr, "\nError:not a excel file input!")
		return
	}

	if *paramProtoOut != "" && *paramPbBinaryOut != "" {
		g.AddOutputType("proto", &converter.ProtoArgs{
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

	glog.Infof("==========%s==========", "开始导出")
	//g.InputFileList

	//fileObjList := make([]*excel.File, 0)

	for _, inputFile := range g.InputFileList {
		file := excel.NewFile(inputFile)
		if file == nil {
			return false
		}

		file.GlobalFD = g.FileDescriptor

		// 电子表格数据导出到Table对象
		if !file.ExportLocalType() {
			return false
		}

		// 整合类型信息和数据
		if !g.AddFile(file) {
			return false
		}

		glog.Infoln(filepath.Base(file.FileName))
		//开始解析数据
		dataModel := model.NewDataModel()

		tab := model.NewTable()
		tab.LocalFD = file.LocalFD
		file.LocalFD.Table = tab

		if !file.ExportData(dataModel) {
			return false
		}

		// 合并所有值到node节点
		if !excel.MergeValues(dataModel, tab, file) {
			return false
		}

	}

	// 根据各种导出类型, 调用各导出器导出
	return g.Convert()
}
