package main

import (
	"fmt"

	"github.com/dave/dst/gendst/data"
	. "github.com/dave/jennifer/jen"
)

const DSTPATH = "github.com/dave/dst"

func generateDecorator(names []string) error {

	f := NewFile("decorator")
	f.ImportName(DSTPATH, "dst")

	f.Func().Params(Id("f").Op("*").Id("fileDecorator")).Id("decorateNode").Params(
		Id("parent").Qual("go/ast", "Node"),
		Id("n").Qual("go/ast", "Node"),
	).Qual(DSTPATH, "Node").BlockFunc(func(g *Group) {
		g.If(List(Id("dn"), Id("ok")).Op(":=").Id("f").Dot("Dst").Dot("Nodes").Index(Id("n")), Id("ok")).Block(
			Return(Id("dn")),
		)
		g.Switch(Id("n").Op(":=").Id("n").Assert(Id("type"))).BlockFunc(func(g *Group) {
			for _, nodeName := range names {
				g.Case(Op("*").Qual("go/ast", nodeName)).BlockFunc(func(g *Group) {

					g.Id("out").Op(":=").Op("&").Qual(DSTPATH, nodeName).Values()

					g.Id("f").Dot("Dst").Dot("Nodes").Index(Id("n")).Op("=").Id("out")
					g.Id("f").Dot("Ast").Dot("Nodes").Index(Id("out")).Op("=").Id("n")

					if nodeName != "Package" {
						g.Line()
						g.Id("out").Dot("Decs").Dot("Before").Op("=").Id("f").Dot("before").Index(Id("n"))
						g.Id("out").Dot("Decs").Dot("After").Op("=").Id("f").Dot("after").Index(Id("n"))
					}
					for _, frag := range data.Info[nodeName] {
						switch frag := frag.(type) {
						case data.Init:
							g.Line().Commentf("Init: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Op("&").Qual(DSTPATH, frag.Type.Name).Values()
						case data.Decoration:
							// nothing here
						case data.String:
							g.Line().Commentf("String: %s", frag.Name)
							if frag.ValueField != nil {
								g.Add(frag.ValueField.Get("out")).Op("=").Add(frag.ValueField.Get("n"))
							}
						case data.Token:
							g.Line().Commentf("Token: %s", frag.Name)
							if frag.TokenField != nil {
								g.Add(frag.TokenField.Get("out")).Op("=").Add(frag.TokenField.Get("n"))
							}
							if frag.ExistsField != nil {
								g.Add(frag.ExistsField.Get("out")).Op("=").Add(frag.Exists.Get("n", true))
							}
						case data.List:
							g.Line().Commentf("List: %s", frag.Name)
							/*
								for _, v := range n.<name> {
									out.<name> = append(out.<name>, d.decorateNode(n, v).(<type>))
								}
							*/
							g.For(List(Id("_"), Id("v")).Op(":=").Range().Add(frag.Field.Get("n"))).Block(
								frag.Field.Get("out").Op("=").Append(
									frag.Field.Get("out"),
									Id("f").Dot("decorateNode").Call(Id("n"), Id("v")).Assert(frag.Elem.Literal(DSTPATH)),
								),
							)
						case data.Map:
							g.Line().Commentf("Map: %s", frag.Name)
							/*
								out.<name> = map[string]<type>{}
								for k, v := range n.<name> {
									out.<name>[k] = d.decorateNode(n, v).(<type>)
								}

								or:

								out.<name> = map[string]<type>{}
								for k, v := range n.<name> {
									out.<name>[k] = d.DecorateObject(v)
								}
							*/
							g.Add(frag.Field.Get("out")).Op("=").Map(String()).Add(frag.Elem.Literal(DSTPATH)).Values()
							g.For(List(Id("k"), Id("v")).Op(":=").Range().Add(frag.Field.Get("n"))).BlockFunc(func(g *Group) {
								if frag.Elem.Name == "Object" {
									// Special case for Package.Imports
									g.Add(frag.Field.Get("out")).Index(Id("k")).Op("=").Id("f").Dot("decorateObject").Call(Id("v"))
								} else {
									g.Add(frag.Field.Get("out")).Index(Id("k")).Op("=").Id("f").Dot("decorateNode").Call(Id("n"), Id("v")).Assert(frag.Elem.Literal(DSTPATH))
								}
							})
						case data.Node:
							g.Line().Commentf("Node: %s", frag.Name)
							/*
								if n.<name> != nil {
									out.<name> = d.decorateNode(n, n.<name>).(<type>)
								}
							*/
							g.If(frag.Field.Get("n").Op("!=").Nil()).Block(
								frag.Field.Get("out").Op("=").Id("f").Dot("decorateNode").Call(Id("n"), frag.Field.Get("n")).Assert(frag.Type.Literal(DSTPATH)),
							)
						case data.Bad:
							g.Line().Comment("Bad")
							g.Add(frag.LengthField.Get("out")).Op("=").Add(frag.Length.Get("n", true))
						case data.Value:
							g.Line().Commentf("Value: %s", frag.Name)
							if frag.Value != nil {
								g.Add(frag.Field.Get("out")).Op("=").Add(frag.Value.Get("n", true))
							} else {
								g.Add(frag.Field.Get("out")).Op("=").Add(frag.Field.Get("n"))
							}
						case data.Scope:
							g.Line().Commentf("Scope: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Id("f").Dot("decorateScope").Call(frag.Field.Get("n"))
						case data.Object:
							g.Line().Commentf("Object: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Id("f").Dot("decorateObject").Call(frag.Field.Get("n"))
						case data.PathDecoration:
							g.Line().Commentf("Path: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Id("f").Dot("resolvePath").Call(Id("parent"), Id("n"))
						case data.SpecialDecoration:
							// ignore
						default:
							panic(fmt.Sprintf("unknown fragment type %T", frag))
						}
					}

					g.Line()
					var found bool
					decs := If(List(Id("nd"), Id("ok")).Op(":=").Id("f").Dot("decorations").Index(Id("n")), Id("ok")).BlockFunc(func(g *Group) {
						for _, frag := range data.Info[nodeName] {
							switch frag := frag.(type) {
							case data.Decoration:
								found = true
								g.If(List(Id("decs"), Id("ok")).Op(":=").Id("nd").Index(Lit(frag.Name)), Id("ok")).Block(
									Id("out").Dot("Decs").Dot(frag.Name).Op("=").Id("decs"),
								)
							}
						}
					})
					if found {
						g.Add(decs)
					}

					g.Line()
					g.Return(Id("out"))

				})
			}
		})
		g.Return(Nil())

	})

	return f.Save("./decorator/decorator-node-generated.go")
}

func generateDecoratorTestHelper(names []string) error {
	f := NewFile("decorator")
	f.ImportName(DSTPATH, "dst")
	f.Func().Id("getDecorationInfo").Params(Id("n").Qual(DSTPATH, "Node")).Params(Id("before"), Id("after").Qual(DSTPATH, "SpaceType"), Id("info").Index().Id("decorationInfo")).BlockFunc(func(g *Group) {
		g.Switch(Id("n").Op(":=").Id("n").Assert(Id("type"))).BlockFunc(func(g *Group) {
			for _, nodeName := range names {
				g.Case(Op("*").Qual(DSTPATH, nodeName)).BlockFunc(func(g *Group) {
					if nodeName != "Package" {
						g.Id("before").Op("=").Id("n").Dot("Decs").Dot("Before")
						g.Id("after").Op("=").Id("n").Dot("Decs").Dot("After")
					}
					for _, frag := range data.Info[nodeName] {
						switch frag := frag.(type) {
						case data.Decoration:
							g.Id("info").Op("=").Append(Id("info"), Id("decorationInfo").Values(Lit(frag.Name), Id("n").Dot("Decs").Dot(frag.Name)))
						}
					}
				})
			}
		})
		g.Return()
	})
	return f.Save("./decorator/decorator-info-generated.go")
}
