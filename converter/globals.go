package converter

import (
	"github.com/821869798/excelconvert/model"
)

type Globals struct {
	ProtoVersion          int
	InputFileList         string
	Converters            []*ConverterContext
	*model.FileDescriptor                         // 类型信息.用于添加各种导出结构
	tableByName           map[string]*model.Table //  防止table重名
	Tables                []*model.Table          // 数据信息.表格数据
}

func NewGlobals() *Globals {
	self := &Globals{
		tableByName:    make(map[string]*model.Table),
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
