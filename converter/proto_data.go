package converter

import (
	"fmt"
	"github.com/821869798/excelconvert/model"
	"github.com/gogo/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"strconv"
)

func buildOneBytesFile(tab *model.Table, PbBinaryOutPath string, pfd *desc.FileDescriptor) bool {
	localFD := tab.LocalFD
	//List数据的结构
	gourpMsg := pfd.FindMessage(fmt.Sprintf("%s.%sTalbe", localFD.Package, localFD.Name))
	gourpDM := dynamic.NewMessage(gourpMsg)
	//单条记录的结构
	recordMsg := pfd.FindMessage(fmt.Sprintf("%s.%s", localFD.Package, localFD.Name))

	for i, r := range tab.Recs {
		recordDM := dynamic.NewMessage(recordMsg)
		for _, node := range r.Nodes {
			if node.Type != model.FieldType_Struct {
				if node.IsRepeated {
					for index, valueNode := range node.Child {
						recordDM.SetRepeatedFieldByName(node.Name, index, getBuildPBValue(node.Type, valueNode))
					}
				} else {
					recordDM.SetFieldByName(node.Name, getBuildPBValue(node.Type, node))
				}
			} else {
				return false
			}
		}
		gourpDM.SetRepeatedFieldByNumber(1, i, recordDM)
	}
	return true
}

func getBuildPBValue(ft model.FieldType, value *model.Node) interface{} {
	switch ft {
	case model.FieldType_Int32:
		v, _ := strconv.ParseInt(value.Value, 10, 32)
		return proto.Int32(int32(v))
	case model.FieldType_UInt32:
		v, _ := strconv.ParseUint(value.Value, 10, 32)
		return proto.Uint32(uint32(v))
	case model.FieldType_Int64:
		v, _ := strconv.ParseInt(value.Value, 10, 64)
		return proto.Int64(v)
	case model.FieldType_UInt64:
		v, _ := strconv.ParseUint(value.Value, 10, 64)
		return proto.Uint64(v)
	case model.FieldType_Float:
		v, _ := strconv.ParseFloat(value.Value, 32)
		return proto.Float32(float32(v))
	case model.FieldType_Bool:
		v, _ := strconv.ParseBool(value.Value)
		return proto.Bool(v)
	case model.FieldType_String:
		return proto.String(value.Value)
	case model.FieldType_Enum:
		return proto.Int32(value.EnumValue)
	default:
		panic("unsupport type" + model.FieldTypeToString(ft))
	}
}
