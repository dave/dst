package main

import (
	"fmt"

	"github.com/dave/dst/gendst/fragment"
	. "github.com/dave/jennifer/jen"
)

func generateFragger(names []string) error {
	f := NewFile("decorator")
	f.Func().Params(Id("f").Op("*").Id("Fragger")).Id("ProcessNode").Params(Id("n").Qual("go/ast", "Node")).Block(
		If(Id("n").Dot("Pos").Call().Dot("IsValid").Call()).Block(
			Id("f").Dot("cursor").Op("=").Int().Parens(Id("n").Dot("Pos").Call()),
		),
		Switch(Id("n").Op(":=").Id("n").Assert(Type())).BlockFunc(func(g *Group) {
			for _, nodeName := range names {
				g.Case(Op("*").Qual("go/ast", nodeName)).BlockFunc(func(g *Group) {
					for _, frag := range fragment.Info[nodeName] {
						switch frag := frag.(type) {
						case fragment.Decoration:
							g.Line().Commentf("Decoration: %s", frag.Name)

							pos := Qual("go/token", "NoPos")
							switch frag.Name {
							case "Start":
								pos = Id("n").Dot("Pos").Call()
							case "End":
								pos = Id("n").Dot("End").Call()
							}

							process := Id("f").Dot("AddDecoration").Call(Id("n"), Lit(frag.Name), pos)

							if frag.Use != nil {
								g.If(frag.Use.Get("n", true)).Block(process)
							} else {
								g.Add(process)
							}
						case fragment.Node:
							g.Line().Commentf("Node: %s", frag.Name)
							g.If(frag.Field.Get("n").Op("!=").Nil()).Block(
								Id("f").Dot("ProcessNode").Call(frag.Field.Get("n")),
							)
						case fragment.List:
							g.Line().Commentf("List: %s", frag.Name)
							g.For(List(Id("_"), Id("v")).Op(":=").Range().Add(frag.Field.Get("n"))).Block(
								Id("f").Dot("ProcessNode").Call(Id("v")),
							)
						case fragment.Map:
							g.Line().Commentf("Map: %s", frag.Name)
							if frag.Elem.Name != "Object" {
								g.For(List(Id("_"), Id("v")).Op(":=").Range().Add(frag.Field.Get("n"))).Block(
									Id("f").Dot("ProcessNode").Call(Id("v")),
								)
							}
						case fragment.Token:
							g.Line().Commentf("Token: %s", frag.Name)
							pos := Qual("go/token", "NoPos")
							if frag.PositionField != nil {
								pos = frag.PositionField.Get("n")
							}
							process := Id("f").Dot("AddToken").Call(Id("n"), frag.Token.Get("n", true), pos)
							if frag.Exists != nil {
								g.If(frag.Exists.Get("n", true)).Block(process)
							} else {
								g.Add(process)
							}
						case fragment.String:
							g.Line().Commentf("String: %s", frag.Name)
							pos := Qual("go/token", "NoPos")
							if frag.PositionField != nil {
								pos = frag.PositionField.Get("n")
							}
							g.Id("f").Dot("AddString").Call(Id("n"), frag.ValueField.Get("n"), pos)
						case fragment.Ignored, fragment.Init, fragment.Value, fragment.Scope, fragment.Object:
							// do nothing
						default:
							panic(fmt.Sprintf("unknown fragment type %T", frag))
						}
					}
					g.Line()
				})
			}
		}),
	)
	return f.Save("./decorator/fragger-generated.go")
}
