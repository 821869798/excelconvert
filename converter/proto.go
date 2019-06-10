package converter

import (
	"bytes"
	"github.com/821869798/excelconvert/model"
	"github.com/golang/glog"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/desc/protoprint"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ProtoArgs struct {
	ProtoOut    string
	PbBinaryOut string
}

type protoFieldDescriptor struct {
	*model.FieldDescriptor

	d *protoDescriptor

	Number int
}

func (self protoFieldDescriptor) Label() string {
	if self.IsRepeated {
		return "repeated "
	}

	if self.d.file.ProtoVersion == 2 {
		return "optional "
	}

	return ""
}

type protoDescriptor struct {
	*model.Descriptor

	ProtoFields []protoFieldDescriptor

	file *protoFileModel
}

type protoFileModel struct {
	Package      string
	ProtoVersion int
	ToolVersion  string
	Messages     []protoDescriptor
	Enums        []protoDescriptor
}

type protoConverter struct {
}

func (self *protoConverter) Run(g *Globals, args interface{}) bool {

	protpArgs := args.(*ProtoArgs)

	var m protoFileModel
	pr := &protoprint.Printer{}

	m.Package = g.FileDescriptor.Pragma.GetString("Package")
	m.ProtoVersion = g.ProtoVersion
	for _, d := range g.FileDescriptor.Descriptors {

		fileName := filepath.Join(protpArgs.ProtoOut, d.Name+".proto")
		file := builder.NewFile(fileName).SetPackageName(m.Package).SetProto3(m.ProtoVersion == 3)

		switch d.Kind {
		case model.DescriptorKind_Struct:
			msg := builder.NewMessage(d.Name)
			for _, fd := range d.Fields {
				fType := getProtoBuildType(fd.TypeString())
				if fType == nil {
					continue
				}
				buildFiled := builder.NewField(fd.Name, fType)
				if fd.IsRepeated {
					buildFiled.SetRepeated()
				}
				msg.AddField(buildFiled)
			}
			file.AddMessage(msg)
			//case model.DescriptorKind_Enum:
			//	m.Enums = append(m.Enums, protoD)
		}

		fd, err := file.Build()
		if err != nil {
			glog.Errorf("%s%s", "导出Proto结构报错:", fileName)
			return false
		}
		var buf bytes.Buffer
		err = pr.PrintProtoFile(fd, &buf)
		if err != nil {
			glog.Errorf("%s%s", "导出Proto文件报错:", fileName)
			return false
		}
		parentPath := filepath.Dir(fileName)
		_ = os.MkdirAll(parentPath, os.ModePerm)
		err = ioutil.WriteFile(fileName, buf.Bytes(), 0777)
		if err != nil {
			glog.Errorln("%s%s", "写入Proto文件报错:", fileName)
		}

	}

	//for _, d := range g.FileDescriptor.Descriptors {
	//	// 这给被限制输出
	//	if !d.File.MatchTag(".proto") {
	//		glog.Infof("%s: %s", "输出器: @Types的'OutputTag'忽略了目标", d.Name)
	//		continue
	//	}
	//
	//	var protoD protoDescriptor
	//	protoD.Descriptor = d
	//	protoD.file = &m
	//
	//	// 遍历字段
	//	for index, fd := range d.Fields {
	//		// 对CombineStruct的XXDefine对应的字段
	//		if d.Usage == model.DescriptorUsage_CombineStruct {
	//
	//			// 这个字段被限制输出
	//			if fd.Complex != nil && !fd.Complex.File.MatchTag(".proto") {
	//				continue
	//			}
	//		}
	//
	//		var field protoFieldDescriptor
	//		field.FieldDescriptor = fd
	//		field.d = &protoD
	//
	//		switch d.Kind {
	//		case model.DescriptorKind_Struct:
	//			field.Number = index + 1
	//		case model.DescriptorKind_Enum:
	//			field.Number = int(fd.EnumValue)
	//		}
	//
	//		protoD.ProtoFields = append(protoD.ProtoFields, field)
	//	}
	//	switch d.Kind {
	//	case model.DescriptorKind_Struct:
	//		m.Messages = append(m.Messages, protoD)
	//	case model.DescriptorKind_Enum:
	//		m.Enums = append(m.Enums, protoD)
	//	}
	//}

	return false
}

func getProtoBuildType(typeString string) *builder.FieldType {
	switch typeString {
	case "int32":
		return builder.FieldTypeInt32()
	case "int64":
		return builder.FieldTypeInt64()
	case "uint32":
		return builder.FieldTypeUInt32()
	case "uint64":
		return builder.FieldTypeUInt64()
	case "float":
		return builder.FieldTypeFloat()
	case "string":
		return builder.FieldTypeString()
	case "bool":
		return builder.FieldTypeBool()
	default:
		return nil
	}
}

func init() {
	RegisterConverter("proto", &protoConverter{})
}
