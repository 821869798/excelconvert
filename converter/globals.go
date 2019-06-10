package converter

import (
	"github.com/821869798/excelconvert/model"
	"github.com/golang/glog"
	"strings"
)

type TableIndex struct {
	Index *model.FieldDescriptor // 表头里的索引
	Row   *model.FieldDescriptor // 索引的数据
}

type Globals struct {
	ProtoVersion          int
	InputFileList         []string
	Converters            []*ConverterContext
	*model.FileDescriptor                         // 类型信息.用于添加各种导出结构
	tableByName           map[string]*model.Table //  防止table重名
	Tables                []*model.Table          // 数据信息.表格数据

	GlobalIndexes []TableIndex      // 类型信息.全局索引
	CombineStruct *model.Descriptor // 类型信息.Combine结构体
}

func NewGlobals() *Globals {
	self := &Globals{
		tableByName:    make(map[string]*model.Table),
		FileDescriptor: model.NewFileDescriptor(),
		CombineStruct:  model.NewDescriptor(),
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

	// 有表格里描述的包名不一致, 无法合成最终的文件
	if self.Pragma.GetString("Package") == "" {
		self.Pragma.SetString("Package", localFD.Pragma.GetString("Package"))
	} else if self.Pragma.GetString("Package") != localFD.Pragma.GetString("Package") {

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

	// 每个表在结构体里的字段
	rowFD := model.NewFieldDescriptor()
	rowFD.Name = localFD.Name
	rowFD.Type = model.FieldType_Struct
	rowFD.Complex = localFD.RowDescriptor()
	rowFD.IsRepeated = true
	rowFD.Order = int32(len(self.CombineStruct.Fields) + 1)

	// 去掉注释中的回车,避免代码生成错误
	rowFD.Comment = strings.Replace(localFD.Name, "\n", " ", -1)
	self.CombineStruct.Add(rowFD)

	if localFD.RowDescriptor() == nil {
		panic("row field null:" + localFD.Name)
	}

	for _, d := range localFD.Descriptors {

		// 非行类型的, 全部忽略
		if d.Usage != model.DescriptorUsage_RowType {
			continue
		}

		for _, indexFD := range d.Indexes {

			key := TableIndex{
				Row:   rowFD,
				Index: indexFD,
			}

			self.GlobalIndexes = append(self.GlobalIndexes, key)

		}

	}

	return true
}
