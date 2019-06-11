package excel

import (
	"github.com/821869798/excelconvert/model"
	"github.com/golang/glog"
	"strconv"
	"strings"
)

type typeCell struct {
	value string
	col   int
}

// 类型表的数据, 数据读取与使用分开使用, 让类型互相没有依赖
type typeModel struct {
	colData map[string]*typeCell

	done bool

	row int

	fd *model.FieldDescriptor

	rawFieldType string
}

func (self *typeModel) getValue(row string) (string, int) {
	if v, ok := self.colData[row]; ok {
		return v.value, v.col
	}

	return "", -1
}

func newTypeModel() *typeModel {
	return &typeModel{
		colData: make(map[string]*typeCell),
		fd:      model.NewFieldDescriptor(),
	}
}

type typeModelRoot struct {
	pragma string

	models []*typeModel

	unknownModel []*typeModel
	fieldTypeCol int

	Col int
	Row int
}

func (self *typeModelRoot) ParsePragma(localFD *model.FileDescriptor) bool {

	if err := localFD.Pragma.Parse(self.pragma); err != nil {
		glog.Error("%s, '%s'", "类型表: 文件特性解析失败", self.pragma)
		return false
	}

	if localFD.Pragma.GetString("TableName") == "" {
		glog.Error("%s", "类型表: 表名(TableName)为空")
		return false
	}

	if localFD.Pragma.GetString("Package") == "" {
		glog.Error("%s", "类型表: 包名(Package)为空")
		return false
	}

	return true
}

// 解析类型表里的所有类型描述
func (self *typeModelRoot) ParseData(localFD *model.FileDescriptor, globalFD *model.FileDescriptor) bool {

	var td *model.Descriptor

	reservedRowFieldType1 := localFD.Pragma.GetString("TableName")
	reservedRowFieldType2 := reservedRowFieldType1 + "Group"

	// 每一行
	for _, m := range self.models {

		self.Row = m.row

		var rawTypeName string

		rawTypeName, self.Col = m.getValue("ObjectType")

		if rawTypeName == reservedRowFieldType1 || rawTypeName == reservedRowFieldType2 {
			glog.Error("%s '%s'", "数据头: 使用了保留的类型名 例如:表名或者表名+Group", rawTypeName)
			return false
		}

		existType, ok := localFD.DescriptorByName[rawTypeName]

		if ok {

			td = existType

		} else {

			td = model.NewDescriptor()
			td.Name = rawTypeName
			localFD.Add(td)
		}

		// 字段名
		m.fd.Name, self.Col = m.getValue("FieldName")

		// 解析类型
		m.rawFieldType, self.Col = m.getValue("FieldType")
		self.fieldTypeCol = self.Col

		fieldType, isrepeated, complexType, ok := findFieldType(localFD, globalFD, m.rawFieldType)
		if !ok {
			return false
		}

		if fieldType == model.FieldType_None {
			self.unknownModel = append(self.unknownModel, m)
		}

		m.fd.Type = fieldType
		m.fd.Complex = complexType
		m.fd.IsRepeated = isrepeated

		var rawFieldValue string
		// 解析值
		rawFieldValue, self.Col = m.getValue("Value")

		kind, enumvalue, ok := parseFieldValue(rawFieldValue)
		if !ok {
			return false
		}

		if td.Kind == model.DescriptorKind_None {
			td.Kind = kind
			// 一些字段有填值, 一些没填值
		} else if td.Kind != kind {
			glog.Error("%s", "类型表: 类型前后不一致, 由枚举值不一致导致")
			return false
		}

		if td.Kind == model.DescriptorKind_Enum {
			if _, ok := td.FieldByNumber[enumvalue]; ok {
				glog.Error("%s %d", "类型表: 重复的枚举值", enumvalue)
				return false
			}
		}

		m.fd.EnumValue = enumvalue

		m.fd.Comment, self.Col = m.getValue("Comment")

		// 去掉注释中的回车,避免代码生成错误
		m.fd.Comment = strings.Replace(m.fd.Comment, "\n", " ", -1)

		var rawMeta string
		rawMeta, self.Col = m.getValue("Meta")

		if err := m.fd.Meta.Parse(rawMeta); err != nil {
			glog.Error("%s, '%s'", "类型表: 字段特性解析失败", err.Error())
			return false
		}

		// 别名
		var rawAlias string
		rawAlias, self.Col = m.getValue("Alias")
		if self.Col != -1 {
			m.fd.Meta.SetString("Alias", rawAlias)
		}

		// 默认值
		var rawDefault string
		rawDefault, self.Col = m.getValue("Default")
		if self.Col != -1 {
			m.fd.Meta.SetString("Default", rawDefault)
		}

		if td.Add(m.fd) != nil {
			glog.Error("%s '%s'", "类型表: 重复的字段名", m.fd.Name)
			return false
		}

	}

	return true
}

func (self *typeModelRoot) SolveUnknownModel(localFD *model.FileDescriptor, globalFD *model.FileDescriptor) bool {

	for _, m := range self.unknownModel {

		self.Row = m.row
		self.Col = self.fieldTypeCol

		fieldType, isrepeatd, complexType, ok := findFieldType(localFD, globalFD, m.rawFieldType)
		if !ok {
			return false
		}

		// 实在是找不到了, 没辙了
		if fieldType == model.FieldType_None {
			glog.Error("%s, '%s'", "类型表: 未知字段类型", m.rawFieldType)
			return false
		}

		m.fd.Type = fieldType
		m.fd.Complex = complexType
		m.fd.IsRepeated = isrepeatd
	}

	return true
}

func findFieldType(localFD *model.FileDescriptor,
	globalFD *model.FileDescriptor,
	rawFieldType string) (model.FieldType, bool, *model.Descriptor, bool) {

	// 开始在本地symbol中找
	testFD := localFD

	for {

		fieldType, isrepeatd, complexType, ok := findlocalFieldType(testFD, rawFieldType)

		if !ok {
			return model.FieldType_None, isrepeatd, nil, false
		}

		if fieldType != model.FieldType_None {
			return fieldType, isrepeatd, complexType, true
		}

		// 已经是全局范围, 说明找不到
		if testFD == globalFD {

			return model.FieldType_None, isrepeatd, nil, true
		}

		// 找不到就换全局范围找
		testFD = globalFD
	}

}

// bool表示是否有错
func findlocalFieldType(localFD *model.FileDescriptor, rawFieldType string) (model.FieldType, bool, *model.Descriptor, bool) {

	var isrepeated bool
	var puretype string

	if strings.HasPrefix(rawFieldType, model.RepeatedKeyword) {

		puretype = rawFieldType[model.RepeatedKeywordLen+1:]

		isrepeated = true
	} else {
		puretype = rawFieldType
	}

	// 解析普通类型
	if ft, ok := model.ParseFieldType(puretype); ok {

		return ft, isrepeated, nil, true

	}

	// 解析内建类型
	if desc, ok := localFD.DescriptorByName[rawFieldType]; ok {

		// 只有枚举( 结构体不允许再次嵌套, 增加理解复杂度 )
		if desc.Kind != model.DescriptorKind_Enum {
			glog.Error("%s, '%s'", "类型表: 结构体字段类型不能是结构体类型", rawFieldType)

			return model.FieldType_None, isrepeated, nil, false
		}

		return model.FieldType_Enum, isrepeated, desc, true

	}

	// 没找到类型, 待二次pass
	return model.FieldType_None, isrepeated, nil, true

}

func parseFieldValue(rawFieldValue string) (model.DescriptorKind, int32, bool) {

	// 非空值是枚举
	if rawFieldValue != "" {

		v, err := strconv.Atoi(rawFieldValue)
		// 解析枚举值
		if err != nil {

			glog.Error("%s, %s", "类型表: 枚举值解析失败", err.Error())
			return model.DescriptorKind_None, 0, false
		}

		return model.DescriptorKind_Enum, int32(v), true
	}

	return model.DescriptorKind_Struct, 0, true
}
