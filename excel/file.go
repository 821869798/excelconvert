package excel

import (
	"github.com/821869798/excelconvert/model"
	"github.com/golang/glog"
	"github.com/tealeg/xlsx"
	"strings"
)

type File struct {
	FileName string
	LocalFD  *model.FileDescriptor // 本文件的类型描述表
	GlobalFD *model.FileDescriptor // 全局的类型描述表
	coreFile *xlsx.File
}

func NewFile(filename string) *File {
	self := &File{
		FileName: filename,
	}

	var err error
	self.coreFile, err = xlsx.OpenFile(filename)
	if err != nil {
		//glog.Error(err.Error())
		glog.Error("%s:%s,%v", "打开excel文件失败", filename, err.Error())
		return nil
	}

	return self
}

func (self *File) ExportLocalType() bool {

	var sheetCount int

	var typeSheet *TypeSheet
	// 解析类型表
	for _, rawSheet := range self.coreFile.Sheets {
		if isTypeSheet(rawSheet.Name) {
			if sheetCount > 0 {
				glog.Error("文件: 类型表在一个表中只能有一份")
				return false
			}

			typeSheet = newTypeSheet(NewSheet(self, rawSheet))

			// 从cell添加类型
			if !typeSheet.Parse(self.LocalFD, self.GlobalFD) {
				return false
			}

			sheetCount++
		}
	}

	if typeSheet == nil {
		glog.Error("%s", "文件: 类型表(@Types)没有找到")
		return false
	}
	return false
}

func isTypeSheet(name string) bool {
	return strings.TrimSpace(name) == model.TypeSheetName
}
