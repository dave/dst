package main

import (
	"fmt"

	"github.com/dave/dst/gendst/data"
	. "github.com/dave/jennifer/jen"
)

// notest

func generateRestorer(names []string) error {

	f := NewFile("decorator")
	f.ImportName(DSTPATH, "dst")
	// func (r *restorer) restoreNode(n dst.Node, allowDuplicate bool) ast.Node {
	// 	switch n := n.(type) {
	// 	case <type>:
	// 		...
	// 	default:
	// 		panic(...)
	// 	}
	// }
	f.Func().Params(Id("r").Op("*").Id("FileRestorer")).Id("restoreNode").Params(
		Id("n").Qual(DSTPATH, "Node"),
		Id("parentName"),
		Id("parentField"),
		Id("parentFieldType").String(),
		Id("allowDuplicate").Bool(),
	).Qual("go/ast", "Node").BlockFunc(func(g *Group) {
		g.If(List(Id("an"), Id("ok")).Op(":=").Id("r").Dot("Ast").Dot("Nodes").Index(Id("n")), Id("ok")).Block(
			If(Id("allowDuplicate")).Block(
				Return(Id("an")),
			).Else().Block(
				Panic(Qual("fmt", "Sprintf").Call(Lit("duplicate node: %#v"), Id("n"))),
			),
		)
		g.Switch(Id("n").Op(":=").Id("n").Assert(Id("type"))).BlockFunc(func(g *Group) {
			for _, nodeName := range names {
				g.Case(Op("*").Qual(DSTPATH, nodeName)).BlockFunc(func(g *Group) {
					if nodeName == "Ident" {
						g.Line()
						g.Comment("Special case for *dst.Ident - replace with SelectorExpr if needed")
						g.Id("sel").Op(":=").Id("r").Dot("restoreIdent").Call(Id("n"), Id("parentName"), Id("parentField"), Id("parentFieldType"), Id("allowDuplicate"))
						g.If(Id("sel").Op("!=").Nil()).Block(
							Return(Id("sel")),
						)
						g.Line()
					}
					g.Id("out").Op(":=").Op("&").Qual("go/ast", nodeName).Values()
					g.Id("r").Dot("Ast").Dot("Nodes").Index(Id("n")).Op("=").Id("out")
					g.Id("r").Dot("Dst").Dot("Nodes").Index(Id("out")).Op("=").Id("n")

					if nodeName != "Package" {
						g.Id("r").Dot("applySpace").Call(Id("n"), Lit("Before"), Id("n").Dot("Decs").Dot("Before"))
					}

					for _, frag := range data.Info[nodeName] {
						switch frag := frag.(type) {
						case data.Init:
							g.Line().Commentf("Init: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Op("&").Qual("go/ast", frag.Type.TypeName()).Values()
						case data.Decoration:
							g.Line().Commentf("Decoration: %s", frag.Name)
							g.Id("r").Dot("applyDecorations").Call(Id("out"), Id("n").Dot("Decs").Dot(frag.Name), Do(func(s *Statement) { s.Lit(frag.Name == "End") }))
						case data.SpecialDecoration:
							g.Line().Commentf("Special decoration: %s", frag.Name)
							g.Id("r").Dot("applyDecorations").Call(Id("out"), frag.Decs.Get("n").Dot(frag.Name), Lit(frag.End))
						case data.Token:
							g.Line().Commentf("Token: %s", frag.Name)
							position := Null()
							value := Null()
							if frag.PositionField != nil {
								if frag.NoPosField != nil {
									position = If(frag.NoPosField.Get("n")).Block(
										frag.PositionField.Get("out").Op("=").Qual("go/token", "NoPos"),
									).Else().Block(
										frag.PositionField.Get("out").Op("=").Id("r").Dot("cursor"),
									)
								} else {
									position = frag.PositionField.Get("out").Op("=").Id("r").Dot("cursor")
								}
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
						case data.String:
							g.Line().Commentf("String: %s", frag.Name)
							if frag.Literal {
								g.Id("r").Dot("applyLiteral").Call(frag.ValueField.Get("n"))
							}
							if frag.PositionField != nil {
								g.Add(frag.PositionField.Get("out")).Op("=").Id("r").Dot("cursor")
							}
							g.Add(frag.ValueField.Get("out")).Op("=").Add(frag.ValueField.Get("n"))
							g.Id("r").Dot("cursor").Op("+=").Qual("go/token", "Pos").Parens(
								Len(frag.ValueField.Get("n")),
							)
						case data.Node:
							g.Line().Commentf("Node: %s", frag.Name)
							/*
								if n.Elt != nil {
									out.Elt = r.restoreNode(n.Elt).(ast.Expr)
								}
							*/
							g.If(frag.Field.Get("n").Op("!=").Nil()).Block(
								frag.Field.Get("out").Op("=").Id("r").Dot("restoreNode").Call(frag.Field.Get("n"), Lit(nodeName), Lit(frag.Field.FieldName()), Lit(frag.Type.TypeName()), Id("allowDuplicate")).Assert(frag.Type.Literal("go/ast")),
							)
						case data.List:
							if frag.NoRestore {
								continue
							}
							g.Line().Commentf("List: %s", frag.Name)
							g.For(List(Id("_"), Id("v")).Op(":=").Range().Add(frag.Field.Get("n"))).Block(
								frag.Field.Get("out").Op("=").Append(
									frag.Field.Get("out"),
									Id("r").Dot("restoreNode").Call(Id("v"), Lit(nodeName), Lit(frag.Field.FieldName()), Lit(frag.Elem.TypeName()), Id("allowDuplicate")).Assert(frag.Elem.Literal("go/ast")),
								),
							)
						case data.Map:
							g.Line().Commentf("Map: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Map(String()).Add(frag.Elem.Literal("go/ast")).Values()
							g.For(List(Id("k"), Id("v")).Op(":=").Range().Add(frag.Field.Get("n"))).BlockFunc(func(g *Group) {
								if frag.Elem.TypeName() == "Object" {
									g.Add(frag.Field.Get("out")).Index(Id("k")).Op("=").Id("r").Dot("restoreObject").Call(Id("v"))
								} else {
									g.Add(frag.Field.Get("out")).Index(Id("k")).Op("=").Id("r").Dot("restoreNode").Call(Id("v"), Lit(nodeName), Lit(frag.Field.FieldName()), Lit(frag.Elem.TypeName()), Id("allowDuplicate")).Assert(frag.Elem.Literal("go/ast"))
								}
							})
						case data.Bad:
							g.Line().Comment("Bad")
							g.Add(frag.FromField.Get("out")).Op("=").Id("r").Dot("cursor")
							g.Id("r").Dot("cursor").Op("+=").Qual("go/token", "Pos").Parens(frag.Length.Get("n", false))
							g.Add(frag.ToField.Get("out")).Op("=").Id("r").Dot("cursor")
						case data.Value:
							g.Line().Commentf("Value: %s", frag.Name)
							if frag.Value != nil {
								g.Add(frag.Field.Get("out")).Op("=").Add(frag.Value.Get("n", false))
							} else {
								g.Add(frag.Field.Get("out")).Op("=").Add(frag.Field.Get("n"))
							}
						case data.Scope:
							g.Line().Commentf("Scope: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Id("r").Dot("restoreScope").Call(frag.Field.Get("n"))
						case data.Object:
							g.Line().Commentf("Object: %s", frag.Name)
							g.Add(frag.Field.Get("out")).Op("=").Id("r").Dot("restoreObject").Call(frag.Field.Get("n"))
						case data.PathDecoration:
							// nothing
						default:
							panic(fmt.Sprintf("unknown fragment type %T", frag))
						}
					}

					if nodeName != "Package" {
						g.Id("r").Dot("applySpace").Call(Id("n"), Lit("After"), Id("n").Dot("Decs").Dot("After"))
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

	return f.Save("./decorator/restorer-generated.go")
}
