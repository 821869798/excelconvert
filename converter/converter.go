package converter

type ConverterContext struct {
	c    Converter
	args interface{}
}

type Converter interface {
	Run(g *Globals)
}

var converterMap = make(map[string]Converter)

func RegisterConverter(ext string, c Converter) {
	if _, ok := converterMap[ext]; ok {
		panic("duplicate coverter")
	}
	converterMap[ext] = c
}
