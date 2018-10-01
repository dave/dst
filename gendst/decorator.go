package main

import (
	"github.com/dave/dst/gendst/fragment"
	. "github.com/dave/jennifer/jen"
)

const DSTPATH = "github.com/dave/dst"

func generateDecorator(names []string) error {

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
			for _, nodeName := range names {
				g.Case(Op("*").Qual("go/ast", nodeName)).BlockFunc(func(g *Group) {
					g.Id("out").Op(":=").Op("&").Qual(DSTPATH, nodeName).Values()
					for _, frag := range fragment.Info[nodeName] {
						switch frag := frag.(type) {
						case fragment.Init:
							g.Line().Commentf("Init: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Op("&").Qual(DSTPATH, frag.Type.Name).Values()
						case fragment.Decoration:
							// nothing here
						case fragment.String:
							g.Line().Commentf("String: %s", frag.Name)
							if frag.ValueField != nil {
								g.Add(frag.ValueField.Get("out")).Op("=").Add(frag.ValueField.Get("n"))
							}
						case fragment.Token:
							g.Line().Commentf("Token: %s", frag.Name)
							if frag.TokenField != nil {
								g.Add(frag.TokenField.Get("out")).Op("=").Add(frag.TokenField.Get("n"))
							}
							if frag.ExistsField != nil {
								g.Add(frag.ExistsField.Get("out")).Op("=").Add(frag.Exists.Get("n", true))
							}
						case fragment.List:
							g.Line().Commentf("List: %s", frag.Name)
							/*
								for _, v := range n.<name> {
									out.<name> = append(out.<name>, d.DecorateNode(v).(<type>))
								}
							*/
							g.For(List(Id("_"), Id("v")).Op(":=").Range().Add(frag.Field.Get("n"))).Block(
								frag.Field.Get("out").Op("=").Append(
									frag.Field.Get("out"),
									Id("d").Dot("DecorateNode").Call(Id("v")).Assert(frag.Elem.Literal(DSTPATH)),
								),
							)
						case fragment.Node:
							g.Line().Commentf("Node: %s", frag.Name)
							/*
								if n.<name> != nil {
									out.<name> = d.DecorateNode(n.<name>).(<type>)
								}
							*/
							g.If(frag.Field.Get("n").Op("!=").Nil()).Block(
								frag.Field.Get("out").Op("=").Id("d").Dot("DecorateNode").Call(frag.Field.Get("n")).Assert(frag.Type.Literal(DSTPATH)),
							)
						case fragment.Ignored:
							// TODO
						case fragment.Value:
							g.Line().Commentf("Value: %s", frag.Name)
							if frag.Value != nil {
								g.Add(frag.Field.Get("out")).Op("=").Add(frag.Value.Get("n", true))
							} else {
								g.Add(frag.Field.Get("out")).Op("=").Add(frag.Field.Get("n"))
							}
						}
					}

					g.Line()
					g.If(List(Id("decs"), Id("ok")).Op(":=").Id("d").Dot("decorations").Index(Id("n")), Id("ok")).Block(
						Id("out").Dot("Decs").Op("=").Id("decs"),
					)

					g.Line()
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
