package main

import (
	"fmt"

	"github.com/dave/dst/gendst/data"
	. "github.com/dave/jennifer/jen"
)

// notest

func generateClone(names []string) error {

	f := NewFilePathName(DSTPATH, "dst")
	f.Comment("Clone returns a deep copy of the node, ready to be re-used elsewhere in the tree.")
	f.Func().Id("Clone").Params(Id("n").Id("Node")).Id("Node").BlockFunc(func(g *Group) {
		g.Switch(Id("n").Op(":=").Id("n").Assert(Id("type"))).BlockFunc(func(g *Group) {
			for _, nodeName := range names {
				g.Case(Op("*").Qual(DSTPATH, nodeName)).BlockFunc(func(g *Group) {
					g.Id("out").Op(":=").Op("&").Id(nodeName).Values()

					if nodeName != "Package" {
						g.Line()
						g.Id("out").Dot("Decs").Dot("Before").Op("=").Id("n").Dot("Decs").Dot("Before")
					}

					for _, frag := range data.Info[nodeName] {
						switch frag := frag.(type) {
						case data.Init:
							g.Line().Commentf("Init: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Op("&").Id(frag.Type.TypeName()).Values()
						case data.Decoration:
							g.Line().Commentf("Decoration: %s", frag.Name)
							g.Id("out").Dot("Decs").Dot(frag.Name).Op("=").Append(Id("out").Dot("Decs").Dot(frag.Name), Id("n").Dot("Decs").Dot(frag.Name).Op("..."))
						case data.Token:
							if frag.NoPosField != nil {
								g.Line().Commentf("Token: %s", frag.Name)
								g.Add(frag.NoPosField.Get("out")).Op("=").Add(frag.NoPosField.Get("n"))
							}
							if frag.TokenField != nil {
								g.Line().Commentf("Token: %s", frag.Name)
								g.Add(frag.TokenField.Get("out")).Op("=").Add(frag.TokenField.Get("n"))
							}
							if frag.ExistsField != nil {
								g.Line().Commentf("Token: %s", frag.Name)
								g.Add(frag.ExistsField.Get("out")).Op("=").Add(frag.ExistsField.Get("n"))
							}
						case data.String:
							g.Line().Commentf("String: %s", frag.Name)
							g.Add(frag.ValueField.Get("out")).Op("=").Add(frag.ValueField.Get("n"))
						case data.Node:
							g.Line().Commentf("Node: %s", frag.Name)
							g.If(frag.Field.Get("n").Op("!=").Nil()).Block(
								frag.Field.Get("out").Op("=").Id("Clone").Call(frag.Field.Get("n")).Assert(frag.Type.Literal(DSTPATH)),
							)
						case data.List:
							g.Line().Commentf("List: %s", frag.Name)
							g.For(List(Id("_"), Id("v")).Op(":=").Range().Add(frag.Field.Get("n"))).Block(
								frag.Field.Get("out").Op("=").Append(
									frag.Field.Get("out"),
									Id("Clone").Call(Id("v")).Assert(frag.Elem.Literal(DSTPATH)),
								),
							)
						case data.Map:
							g.Line().Commentf("Map: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Map(String()).Add(frag.Elem.Literal(DSTPATH)).Values()
							g.For(List(Id("k"), Id("v")).Op(":=").Range().Add(frag.Field.Get("n"))).BlockFunc(func(g *Group) {
								if frag.Elem.TypeName() == "Object" {
									g.Add(frag.Field.Get("out")).Index(Id("k")).Op("=").Id("CloneObject").Call(Id("v"))
								} else {
									g.Add(frag.Field.Get("out")).Index(Id("k")).Op("=").Id("Clone").Call(Id("v")).Assert(frag.Elem.Literal(DSTPATH))
								}
							})
						case data.Value:
							g.Line().Commentf("Value: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Add(frag.Field.Get("n"))
						case data.Scope:
							g.Line().Commentf("Scope: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Id("CloneScope").Call(frag.Field.Get("n"))
						case data.Object:
							g.Line().Commentf("Object: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Id("CloneObject").Call(frag.Field.Get("n"))
						case data.Bad:
							g.Line().Comment("Bad")
							g.Add(frag.LengthField.Get("out")).Op("=").Add(frag.LengthField.Get("n"))
						case data.PathDecoration:
							g.Line().Commentf("Path: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Add(frag.Field.Get("n"))
						case data.SpecialDecoration:
							// ignore
						default:
							panic(fmt.Sprintf("unknown fragment type %T", frag))
						}
					}

					if nodeName != "Package" {
						g.Line()
						g.Id("out").Dot("Decs").Dot("After").Op("=").Id("n").Dot("Decs").Dot("After")
					}

					g.Line()
					g.Return(Id("out"))
				})
			}
			g.Default().Block(
				Panic(Qual("fmt", "Sprintf").Call(Lit("%T"), Id("n"))),
			)
		})
	})

	return f.Save("./clone-generated.go")
}
