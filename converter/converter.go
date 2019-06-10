package converter

import "github.com/golang/glog"

type ConverterContext struct {
	c    Converter
	args interface{}
}

func (self *ConverterContext) Start(g *Globals) bool {
	glog.Infoln(self.args)

	return self.c.Run(g, self.args)
}

type Converter interface {
	Run(g *Globals, args interface{}) bool
}

var converterMap = make(map[string]Converter)

func RegisterConverter(ext string, c Converter) {
	if _, ok := converterMap[ext]; ok {
		panic("duplicate coverter")
	}
	converterMap[ext] = c
}
