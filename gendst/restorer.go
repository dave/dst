package main

import (
	"github.com/dave/dst/gendst/fragment"
	. "github.com/dave/jennifer/jen"
)

func generateRestorer(names []string) error {

	f := NewFile("decorator")
	f.ImportName(DSTPATH, "dst")
	// func (r *restorer) restoreNode(n dst.Node) ast.Node {
	// 	switch n := n.(type) {
	// 	case <type>:
	// 		...
	// 	default:
	// 		panic(...)
	// 	}
	// }
	f.Func().Params(Id("r").Op("*").Id("fileRestorer")).Id("restoreNode").Params(Id("n").Qual(DSTPATH, "Node")).Qual("go/ast", "Node").BlockFunc(func(g *Group) {
		g.If(List(Id("an"), Id("ok")).Op(":=").Id("r").Dot("nodes").Index(Id("n")), Id("ok")).Block(
			Return(Id("an")),
		)
		g.Switch(Id("n").Op(":=").Id("n").Assert(Id("type"))).BlockFunc(func(g *Group) {
			for _, nodeName := range names {
				g.Case(Op("*").Qual(DSTPATH, nodeName)).BlockFunc(func(g *Group) {
					g.Id("out").Op(":=").Op("&").Qual("go/ast", nodeName).Values()
					for _, frag := range fragment.Info[nodeName] {
						switch frag := frag.(type) {
						case fragment.Init:
							g.Line().Commentf("Init: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Op("&").Qual("go/ast", frag.Type.Name).Values()
						case fragment.Decoration:
							g.Line().Commentf("Decoration: %s", frag.Name)
							g.Id("r").Dot("applyDecorations").Call(Id("n").Dot("Decs").Dot(frag.Name))
						case fragment.Token:
							g.Line().Commentf("Token: %s", frag.Name)
							position := Null()
							value := Null()
							if frag.PositionField != nil {
								position = frag.PositionField.Get("out").Op("=").Id("r").Dot("cursor")
							}
							if frag.TokenField != nil {
								value = frag.TokenField.Get("out").Op("=").Add(frag.Token.Get("n", false))
							}
							action := Id("r").Dot("cursor").Op("+=").Qual("go/token", "Pos").Parens(
								Len(frag.Token.Get("n", false).Dot("String").Call()),
							)
							if frag.Exists != nil {
								g.If(frag.Exists.Get("n", false)).Block(value, position, action)
							} else {
								g.Add(value)
								g.Add(position)
								g.Add(action)
							}
						case fragment.String:
							g.Line().Commentf("String: %s", frag.Name)
							if frag.PositionField != nil {
								g.Add(frag.PositionField.Get("out")).Op("=").Id("r").Dot("cursor")
							}
							if frag.ValueField != nil {
								g.Add(frag.ValueField.Get("out")).Op("=").Add(frag.ValueField.Get("n"))
							}
							g.Id("r").Dot("cursor").Op("+=").Qual("go/token", "Pos").Parens(
								Len(frag.ValueField.Get("n")),
							)
						case fragment.Node:
							g.Line().Commentf("Node: %s", frag.Name)
							/*
								if n.Elt != nil {
									out.Elt = r.restoreNode(n.Elt).(ast.Expr)
								}
							*/
							g.If(frag.Field.Get("n").Op("!=").Nil()).Block(
								frag.Field.Get("out").Op("=").Id("r").Dot("restoreNode").Call(frag.Field.Get("n")).Assert(frag.Type.Literal("go/ast")),
							)
						case fragment.List:
							g.Line().Commentf("List: %s", frag.Name)
							g.For(List(Id("_"), Id("v")).Op(":=").Range().Add(frag.Field.Get("n"))).Block(
								frag.Field.Get("out").Op("=").Append(
									frag.Field.Get("out"),
									Id("r").Dot("restoreNode").Call(Id("v")).Assert(frag.Elem.Literal("go/ast")),
								),
							)
						case fragment.Ignored:
							// TODO
						case fragment.Value:
							g.Line().Commentf("Value: %s", frag.Name)
							if frag.Value != nil {
								g.Add(frag.Field.Get("out")).Op("=").Add(frag.Value.Get("n", false))
							} else {
								g.Add(frag.Field.Get("out")).Op("=").Add(frag.Field.Get("n"))
							}
						}
					}
					g.Line()
					g.Return(Id("out"))
				})
			}
			g.Default().Block(
				Panic(Qual("fmt", "Sprintf").Call(Lit("%T"), Id("n"))),
			)
		})
		g.Return(Nil())
	})

	return f.Save("./decorator/restorer-generated.go")
}
