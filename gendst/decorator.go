package main

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
)

const DSTPATH = "github.com/dave/dst"

func generateDecorator(names []string, nodes map[string]NodeInfo) error {

	f := NewFile("decorator")
	f.ImportName(DSTPATH, "dst")

	/*
		func (d *Decorator) DecorateNode(n ast.Node) dst.Node {
			if dn, ok := d.nodes[n]; ok {
				return dn
			}
		}
	*/
	f.Func().Params(Id("d").Op("*").Id("Decorator")).Id("DecorateNode").Params(
		Id("n").Qual("go/ast", "Node"),
	).Qual(DSTPATH, "Node").BlockFunc(func(g *Group) {
		g.If(List(Id("dn"), Id("ok")).Op(":=").Id("d").Dot("nodes").Index(Id("n")), Id("ok")).Block(
			Return(Id("dn")),
		)
		g.Switch(Id("n").Op(":=").Id("n").Assert(Id("type"))).BlockFunc(func(g *Group) {
			for _, name := range names {
				node := nodes[name]
				g.Case(Op("*").Qual("go/ast", name)).BlockFunc(func(g *Group) {
					g.Id("out").Op(":=").Op("&").Qual(DSTPATH, name).Values()

					for _, frag := range node.Fragments {
						switch frag.AstType {
						case "Node":
							g.If(Id("n").Dot(frag.Name).Op("!=").Nil()).Block(
								Id("out").Dot(frag.Name).Op("=").Id("d").Dot("DecorateNode").Call(
									Id("n").Dot(frag.Name),
								).Assert(typeLiteral(DSTPATH, frag.DstTypeActual, frag.DstTypePointer)),
							)
						case "[]Node":
							g.For(List(Id("_"), Id("v")).Op(":=").Range().Id("n").Dot(frag.Name)).Block(
								Id("out").Dot(frag.Name).Op("=").Append(
									Id("out").Dot(frag.Name),
									Id("d").Dot("DecorateNode").Call(Id("v")).Assert(typeLiteral(DSTPATH, frag.DstTypeActual, frag.DstTypePointer)),
								),
							)
						case "string", "bool", "Token":
							g.Id("out").Dot(frag.Name).Op("=").Id("n").Dot(frag.Name)
						case "Pos":
							if frag.DstType == "bool" {
								g.If(Id("n").Dot(frag.Name).Op("!=").Qual("go/token", "NoPos")).Block(
									Id("out").Dot(frag.Name).Op("=").True(),
								)
							}
						default:
							panic(fmt.Sprintf("%s: %s", frag.Name, frag.AstType))
						}
					}
					for _, field := range node.Data {
						switch field.Type {
						case "[]Node":
							g.For(List(Id("_"), Id("v")).Op(":=").Range().Id("n").Dot(field.Name)).Block(
								Id("out").Dot(field.Name).Op("=").Append(
									Id("out").Dot(field.Name),
									Id("d").Dot("DecorateNode").Call(Id("v")).Assert(typeLiteral(DSTPATH, field.Actual, field.Pointer)),
								),
							)
						case "ChanDir":
							g.Id("out").Dot(field.Name).Op("=").Qual(DSTPATH, "ChanDir").Parens(Id("n").Dot(field.Name))
						case "Token", "bool":
							g.Id("out").Dot(field.Name).Op("=").Id("n").Dot(field.Name)
						case "Object":
							// TODO
							g.Commentf("TODO: %s (Object)", field.Name)
						case "Scope":
							// TODO
							g.Commentf("TODO: %s (Scope)", field.Name)
						default:
							panic(fmt.Sprintf("%s: %s", field.Name, field.Type))
						}
					}
					if node.FromToLength {
						g.Id("out").Dot("Length").Op("=").Int().Parens(Id("n").Dot("End").Call().Op("-").Id("n").Dot("Pos").Call())
					}
					g.If(List(Id("decs"), Id("ok")).Op(":=").Id("d").Dot("decorations").Index(Id("n")), Id("ok")).Block(
						Id("out").Dot("Decs").Op("=").Id("decs"),
					)
					g.Id("d").Dot("nodes").Index(Id("n")).Op("=").Id("out")
					g.Return(Id("out"))

				})
			}
		})
		g.Return(Nil())

	})

	return f.Save("./decorator/decorator-generated.go")
}

func typeLiteral(path, actual string, pointer bool) *Statement {
	return Do(func(s *Statement) {
		if pointer {
			s.Op("*")
		}
	}).Qual(path, actual)
}
