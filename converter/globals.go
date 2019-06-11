package converter

import (
	"github.com/821869798/excelconvert/excel"
	"github.com/821869798/excelconvert/model"
	"github.com/golang/glog"
)

type TableIndex struct {
	Index *model.FieldDescriptor // 表头里的索引
	Row   *model.FieldDescriptor // 索引的数据
}

type Globals struct {
	ProtoVersion          int
	InputFileList         []string
	Converters            []*ConverterContext
	*model.FileDescriptor                        //  类型信息.用于添加各种导出结构
	fileByName            map[string]*excel.File //  所有文件信息，防止Table重名
	Files                 []*excel.File          //  所有文件
}

func NewGlobals() *Globals {
	self := &Globals{
		fileByName:     make(map[string]*excel.File),
		FileDescriptor: model.NewFileDescriptor(),
	}
	return self
}

func (self *Globals) AddOutputType(name string, args interface{}) {

	if c, ok := converterMap[name]; ok {
		self.Converters = append(self.Converters, &ConverterContext{
			c:    c,
			args: args,
		})
	} else {
		panic("output type not found:" + name)
	}

}

func (self *Globals) Convert() bool {

	glog.Infof("==========%s==========", "解析完成，开始导出")

	for _, c := range self.Converters {

		if !c.Start(self) {
			return false
		}
	}

	return true

}

func (self *Globals) AddFile(file *excel.File) bool {
	if _, ok := self.fileByName[file.LocalFD.Name]; ok {
		glog.Errorf("%s, '%s'", "合并: 表名(TableName)重复", file.LocalFD.Name)
		return false
	}

	self.fileByName[file.LocalFD.Name] = file
	self.Files = append(self.Files, file)
	return true
}

//func (self *Globals) AddGlobalTypes(localFD *model.FileDescriptor) bool {
//
//	// 有表格里描述的包名不一致, 无法合成最终的文件
//	if self.Package == "" {
//		self.Pragma.SetString("Package", localFD.Pragma.GetString("Package"))
//	} else if self.Package != localFD.Package {
//
//		glog.Errorf("%s, '%s' '%s'", "合并: 所有表中的@Types中的包名(Package)请保持一致", localFD.Pragma.GetString("TableName"), self.Pragma.GetString("TableType"))
//		return false
//	}
//
//	// 将行定义结构也添加到文件中
//	for _, d := range localFD.Descriptors {
//		if !self.FileDescriptor.Add(d) {
//			glog.Errorf("%s, %s", "合并: 重复的类型名(表名?)", d.Name)
//			return false
//		}
//	}
//
//	return true
//}
