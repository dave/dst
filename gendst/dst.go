package main

import (
	. "github.com/dave/jennifer/jen"
)

func generateDst(names []string) error {
	f := NewFile("dst")
	for _, name := range names {
		// func (v *<name>) Decorations() []Decoration {
		// 	return v.Decs
		// }
		f.Func().Params(Id("v").Op("*").Id(name)).Id("Decorations").Params().Index().Id("Decoration").Block(
			Return(Id("v").Dot("Decs")),
		)
	}
	return f.Save("./generated.go")
}
