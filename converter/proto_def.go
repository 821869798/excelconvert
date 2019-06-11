package converter

import (
	"bytes"
	"github.com/821869798/excelconvert/model"
	"github.com/golang/glog"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"github.com/jhump/protoreflect/desc/protoprint"
	"io/ioutil"
	"os"
	"path/filepath"
)

func buildOneProtoFiler(g *Globals, protoOutPath string, localFD *model.FileDescriptor) (bool, *desc.FileDescriptor) {
	pr := &protoprint.Printer{}
	Package := localFD.Package
	ProtoVersion := g.ProtoVersion
	baseName := localFD.Name + ".proto"
	fileName := filepath.Join(protoOutPath, baseName)
	file := builder.NewFile(fileName).SetPackageName(Package).SetProto3(ProtoVersion == 3)
	//复杂类型
	complexMap := make(map[string]interface{})
	for _, d := range localFD.Descriptors {
		switch d.Kind {
		case model.DescriptorKind_Struct:
			msg := builder.NewMessage(d.Name)
			for _, fd := range d.Fields {
				fType := getProtoBuildType(fd.TypeString(), complexMap)
				if fType == nil {
					continue
				}
				buildFiled := builder.NewField(fd.Name, fType)
				if fd.IsRepeated {
					buildFiled.SetRepeated()
				}
				msg.AddField(buildFiled)
			}
			complexMap[d.Name] = msg
			file.AddMessage(msg)
		case model.DescriptorKind_Enum:
			en := builder.NewEnum(d.Name)
			for _, fd := range d.Fields {
				en.AddValue(builder.NewEnumValue(fd.Name))
			}
			complexMap[d.Name] = en
			file.AddEnum(en)
		}
	}

	pfd, err := file.Build()
	if err != nil {
		glog.Errorf("%s%s", "导出Proto结构报错:", localFD.Name)
		return false, nil
	}
	var buf bytes.Buffer
	err = pr.PrintProtoFile(pfd, &buf)
	if err != nil {
		glog.Errorf("%s%s", "导出Proto文件报错:", localFD.Name)
		return false, nil
	}
	parentPath := filepath.Dir(fileName)
	_ = os.MkdirAll(parentPath, os.ModePerm)
	err = ioutil.WriteFile(fileName, buf.Bytes(), 0777)
	if err != nil {
		glog.Errorln("%s%s", "写入Proto文件报错:", localFD.Name)
	}

	return true, pfd
}

func getProtoBuildType(typeString string, complexMap map[string]interface{}) *builder.FieldType {
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
		if complexType, ok := complexMap[typeString]; ok {
			switch complexType.(type) {
			case *builder.MessageBuilder:
				return builder.FieldTypeMessage(complexType.(*builder.MessageBuilder))
			case *builder.EnumBuilder:
				return builder.FieldTypeEnum(complexType.(*builder.EnumBuilder))
			}
		}
		return nil
	}
}
