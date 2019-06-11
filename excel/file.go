package excel

import (
	"github.com/821869798/excelconvert/model"
	"github.com/golang/glog"
	"github.com/tealeg/xlsx"
	"path/filepath"
	"strings"
)

// 检查单元格值重复结构
type valueRepeatData struct {
	fd    *model.FieldDescriptor
	value string
}

type File struct {
	FileName     string                //file完整的路径
	BaseFileName string                //文件名,不包括后缀
	LocalFD      *model.FileDescriptor // 本文件的类型描述表
	GlobalFD     *model.FileDescriptor // 全局的类型描述表
	coreFile     *xlsx.File

	Header *DataHeader
	Data   *DataSheet

	valueRepByKey map[valueRepeatData]bool // 检查单元格值重复map
}

func NewFile(filename string) *File {
	self := &File{
		valueRepByKey: make(map[valueRepeatData]bool),
		LocalFD:       model.NewFileDescriptor(),
		FileName:      filename,
		BaseFileName:  strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename)),
	}

	var err error
	self.coreFile, err = xlsx.OpenFile(filename)
	if err != nil {
		//glog.Error(err.Error())
		glog.Errorf("%s:%s,%v", "打开excel文件失败", filename, err.Error())
		return nil
	}

	return self
}

func (self *File) GlobalFileDesc() *model.FileDescriptor {
	return self.GlobalFD

}

func (self *File) ExportLocalType() bool {

	var sheetCount int

	var typeSheet *TypeSheet
	// 解析类型表
	for _, rawSheet := range self.coreFile.Sheets {
		if isTypeSheet(rawSheet.Name) {
			if sheetCount > 0 {
				glog.Errorf("文件: 类型表在一个表中只能有一份")
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
		glog.Errorf("%s", "文件: 类型表(@Types)没有找到")
		return false
	}

	for _, rawSheet := range self.coreFile.Sheets {
		// 是数据表
		if !isTypeSheet(rawSheet.Name) {
			dSheet := newDataSheet(NewSheet(self, rawSheet))
			if !dSheet.Valid() {
				continue
			}

			glog.Infof("            %s", rawSheet.Name)

			dataHeader := newDataHeadSheet()

			// 检查引导头
			if !dataHeader.ParseProtoField(dSheet.Sheet, self.LocalFD, self.GlobalFD) {
				return false
			}

			self.Header = dataHeader
			self.Data = dSheet

			break
		}
	}

	// File描述符的名字必须放在类型里, 因为这里始终会被调用, 但是如果数据表缺失, 是不会更新Name的
	self.LocalFD.Name = self.LocalFD.Pragma.GetString("TableName")
	self.LocalFD.Package = self.LocalFD.Pragma.GetString("Package")

	return true
}

func (self *File) ExportData(dataModel *model.DataModel) bool {

	glog.Infof("            %s", self.Data.Name)

	if !self.Data.Export(self, dataModel) {
		return false
	}

	return true

}

func (self *File) CheckValueRepeat(fd *model.FieldDescriptor, value string) bool {

	key := valueRepeatData{
		fd:    fd,
		value: value,
	}

	if _, ok := self.valueRepByKey[key]; ok {
		return false
	}

	self.valueRepByKey[key] = true

	return true
}

func isTypeSheet(name string) bool {
	return strings.TrimSpace(name) == model.TypeSheetName
}
