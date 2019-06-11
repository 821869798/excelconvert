package converter

import (
	"github.com/821869798/excelconvert/model"
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
	protoOutPath := protpArgs.ProtoOut
	PbBinaryOutPath := protpArgs.PbBinaryOut

	for _, fd := range g.Files {
		ok, pfd := buildOneProtoFiler(g, protoOutPath, fd)
		if !ok {
			return false
		}
		if !buildOneBytesFile(fd.Table, PbBinaryOutPath, pfd) {
			return false
		}
	}

	return true

	//buildAllProtoFile(g, protpArgs.ProtoOut)
	//
	//buildAllBytesFile(g, protpArgs.PbBinaryOut)

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

func init() {
	RegisterConverter("proto", &protoConverter{})
}
