package converter

type ProtoArgs struct {
	PackageName string
	TableName   string
	ProtoOut    string
	PbBinaryOut string
}

type protoConverter struct {
}

func (self *protoConverter) Run(g *Globals) {

}

func init() {
	RegisterConverter("proto", &protoConverter{})
}
