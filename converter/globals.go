package converter

import (
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
	*model.FileDescriptor                                  //  类型信息.用于添加各种导出结构
	tableByName           map[string]*model.Table          //  防止table重名
	Tables                []*model.Table                   //  数据信息.表格数据
	fileByName            map[string]*model.FileDescriptor //  所有文件信息，防止重名
	Files                 []*model.FileDescriptor          //  所有文件
}

func NewGlobals() *Globals {
	self := &Globals{
		tableByName:    make(map[string]*model.Table),
		fileByName:     make(map[string]*model.FileDescriptor),
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

func (self *Globals) AddTypes(localFD *model.FileDescriptor) bool {
	if _, ok := self.fileByName[localFD.Name]; ok {
		glog.Errorf("%s, '%s'", "合并: 表名(TableName)重复", localFD.Name)
		return false
	}

	self.fileByName[localFD.Name] = localFD
	self.Files = append(self.Files, localFD)
	return true
}

func (self *Globals) AddGlobalTypes(localFD *model.FileDescriptor) bool {

	// 有表格里描述的包名不一致, 无法合成最终的文件
	if self.Package == "" {
		self.Pragma.SetString("Package", localFD.Pragma.GetString("Package"))
	} else if self.Package != localFD.Package {

		glog.Errorf("%s, '%s' '%s'", "合并: 所有表中的@Types中的包名(Package)请保持一致", localFD.Pragma.GetString("TableName"), self.Pragma.GetString("TableType"))
		return false
	}

	// 将行定义结构也添加到文件中
	for _, d := range localFD.Descriptors {
		if !self.FileDescriptor.Add(d) {
			glog.Errorf("%s, %s", "合并: 重复的类型名(表名?)", d.Name)
			return false
		}
	}

	return true
}

// 合并每个表带的类型
func (self *Globals) AddContent(tab *model.Table) bool {

	localFD := tab.LocalFD

	if _, ok := self.tableByName[localFD.Name]; ok {

		glog.Errorf("%s, '%s'", "合并: 表名(TableName)重复", localFD.Name)
		return false
	}

	// 表的全局类型信息与合并信息一致
	tab.GlobalFD = self.FileDescriptor

	self.tableByName[localFD.Name] = tab
	self.Tables = append(self.Tables, tab)

	return true
}
