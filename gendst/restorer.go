package main

import (
	"fmt"
	"go/types"

	. "github.com/dave/jennifer/jen"
)

func generateRestorer(names []string, nodes map[string]NodeInfo) error {

	f := NewFile("decorator")
	f.ImportName(DSTPATH, "dst")
	// func (r *Restorer) RestoreNode(n dst.Node) ast.Node {
	// 	switch n := n.(type) {
	// 	case <type>:
	// 		...
	// 	default:
	// 		panic(...)
	// 	}
	// }
	f.Func().Params(Id("r").Op("*").Id("FileRestorer")).Id("RestoreNode").Params(Id("n").Qual(DSTPATH, "Node")).Qual("go/ast", "Node").BlockFunc(func(g *Group) {
		g.If(List(Id("an"), Id("ok")).Op(":=").Id("r").Dot("nodes").Index(Id("n")), Id("ok")).Block(
			Return(Id("an")),
		)
		g.Switch(Id("n").Op(":=").Id("n").Assert(Id("type"))).BlockFunc(func(g *Group) {
			for _, name := range names {
				node := nodes[name]
				g.Case(Op("*").Qual(DSTPATH, name)).BlockFunc(func(g *Group) {

					if name == "FuncDecl" {
						// SPECIAL CASE
						g.Return(Id("r").Dot("funcDeclOverride").Call(Id("n")))
						return
					}

					g.Id("r").Dot("applyDecorations").Call(
						Id("n").Dot("Decs"),
						Lit(""),
						Lit(false),
					)

					g.Id("out").Op(":=").Op("&").Qual("go/ast", name).Values()

					if node.FromToLength {
						g.Id("out").Dot("From").Op("=").Id("r").Dot("cursor")
						g.Id("r").Dot("cursor").Op("+=").Qual("go/token", "Pos").Parens(Id("n").Dot("Length"))
						g.Id("out").Dot("To").Op("=").Id("r").Dot("cursor")
					}

					for _, frag := range node.Fragments {

						g.BlockFunc(func(g *Group) {
							// Apply the Before Fragment decorators
							g.Id("r").Dot("applyDecorations").Call(
								Id("n").Dot("Decs"),
								Lit(frag.Name),
								Lit(false),
							)

							g.List(Id("prefix"), Id("length"), Id("suffix")).Op(":=").Id("getLength").Call(Id("n"), Lit(frag.Name))

							// Record the cursor position if there's a position field for this fragment
							if frag.AstPositionField != "" {
								if frag.DstType == "bool" {
									g.If(Id("n").Dot(frag.Name)).Block(
										Id("out").Dot(frag.AstPositionField).Op("=").Id("r").Dot("cursor"),
									)
								} else {
									g.Id("out").Dot(frag.AstPositionField).Op("=").Id("r").Dot("cursor")
								}
							}

							// Increment the cursor if there's a prefix length
							g.Id("r").Dot("cursor").Op("+=").Qual("go/token", "Pos").Parens(Id("prefix"))

							// Copy the values
							switch frag.AstType {
							case "Node":
								g.If(Id("n").Dot(frag.Name).Op("!=").Nil()).Block(
									Id("out").Dot(frag.Name).Op("=").Id("r").Dot("RestoreNode").Call(Id("n").Dot(frag.Name)).Assert(typeLiteral("go/ast", frag.AstTypeActual, frag.AstTypePointer)),
								)
							case "[]Node":
								g.For(List(Id("_"), Id("v")).Op(":=").Range().Id("n").Dot(frag.Name)).Block(
									Id("out").Dot(frag.Name).Op("=").Append(
										Id("out").Dot(frag.Name),
										Id("r").Dot("RestoreNode").Call(Id("v")).Assert(typeLiteral("go/ast", frag.AstTypeActual, frag.AstTypePointer)),
									),
								)
							case "Token", "string", "bool":
								g.Id("out").Dot(frag.Name).Op("=").Id("n").Dot(frag.Name)
							}

							// Increment the cursor if there's a fixed / Token / string length
							g.Id("r").Dot("cursor").Op("+=").Qual("go/token", "Pos").Parens(Id("length"))

							// Increment the cursor if there's a suffix length
							g.Id("r").Dot("cursor").Op("+=").Qual("go/token", "Pos").Parens(Id("suffix"))

							// Apply the After Fragment decorators
							g.Id("r").Dot("applyDecorations").Call(
								Id("n").Dot("Decs"),
								Lit(frag.Name),
								Lit(true),
							)
						})

					}

					for _, field := range node.Data {
						switch field.Type {
						case "[]Node":
							g.For(List(Id("_"), Id("v")).Op(":=").Range().Id("n").Dot(field.Name)).Block(
								Id("out").Dot(field.Name).Op("=").Append(
									Id("out").Dot(field.Name),
									Id("r").Dot("RestoreNode").Call(Id("v")).Assert(typeLiteral("go/ast", field.Actual, field.Pointer)),
								),
							)
						case "ChanDir":
							g.Id("out").Dot(field.Name).Op("=").Qual("go/ast", "ChanDir").Parens(Id("n").Dot(field.Name))
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

					g.Id("r").Dot("applyDecorations").Call(
						Id("n").Dot("Decs"),
						Lit(""),
						Lit(true),
					)

					if name == "CommentGroup" {
						g.Id("r").Dot("Comments").Op("=").Append(Id("r").Dot("Comments"), Id("out"))
					}

					g.Id("r").Dot("nodes").Index(Id("n")).Op("=").Id("out")

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

func getTypeName(t types.Type, node *types.Interface) string {
	var slice string
	if sl, ok := t.(*types.Slice); ok {
		slice = "[]"
		t = sl.Elem()
	}
	if types.Implements(t, node) {
		return slice + "Node"
	}
	t = unwrap(t)
	var name string
	switch t := t.(type) {
	case *types.Named:
		name = t.Obj().Name()
	case *types.Basic:
		name = t.String()
	}
	return slice + name
}

func getTypeActual(t types.Type) (actual string, pointer bool) {
	if sl, ok := t.(*types.Slice); ok {
		t = sl.Elem()
	}
	if ptr, ok := t.(*types.Pointer); ok {
		pointer = true
		t = ptr.Elem()
	}
	switch t := t.(type) {
	case *types.Named:
		actual = t.Obj().Name()
	case *types.Basic:
		actual = t.String()
	}
	return
}

// This was used to generate the basis of restorer-length.go, but special cases are added after
// generation, so we should not run it again.
func generateLength(names []string, nodes map[string]NodeInfo) error {
	f := NewFile("decorator")
	f.ImportName(DSTPATH, "dst")

	/*
		func getLength(node dst.Node, fragment string) (suffix, length, prefix int) {
			switch node := node.(type) {
			case *dst.<name>:
				switch fragment {
				case <frag.Name>:
					return ...
				}
			}
		}
	*/
	f.Func().Id("getLength").Params(Id("n").Qual(DSTPATH, "Node"), Id("fragment").String()).Params(List(Id("suffix"), Id("length"), Id("prefix")).Int()).BlockFunc(func(g *Group) {
		g.Switch(Id("n").Op(":=").Id("n").Assert(Type())).BlockFunc(func(g *Group) {
			for _, name := range names {
				node := nodes[name]
				g.Case(Op("*").Qual(DSTPATH, name)).Block(
					Switch(Id("fragment")).BlockFunc(func(g *Group) {
						for _, frag := range node.Fragments {
							prefix, _ := matchInt(prefixLengths, name, frag.Name)
							suffix, _ := matchInt(suffixLengths, name, frag.Name)
							length, hasFixedLength := matchInt(fixedLength, name, frag.Name)
							var lengthStatement *Statement
							var lengthIsZero bool
							if hasFixedLength {
								lengthIsZero = length == 0
								lengthStatement = Lit(length)
							} else if frag.AstType == "Token" {
								lengthStatement = Len(Id("n").Dot(frag.Name).Dot("String").Call())
							} else if frag.AstType == "string" {
								lengthStatement = Len(Id("n").Dot(frag.Name))
							} else {
								lengthIsZero = true
								lengthStatement = Lit(0)
							}
							allValuesZero := prefix == 0 && suffix == 0 && lengthIsZero
							g.Case(Lit(frag.Name)).BlockFunc(func(g *Group) {
								retVal := Return(Lit(prefix), lengthStatement, Lit(suffix))
								retZero := Return(Lit(0), Lit(0), Lit(0))
								if allValuesZero {
									g.Add(retVal)
								} else {
									switch frag.DstType {
									case "Node":
										g.If(Id("n").Dot(frag.Name).Op("!=").Nil()).Block(retVal)
										g.Add(retZero)
									case "[]Node":
										g.If(Len(Id("n").Dot(frag.Name)).Op(">").Lit(0)).Block(retVal)
										g.Add(retZero)
									case "Token":
										g.If(Id("n").Dot(frag.Name).Op("!=").Qual("go/token", "ILLEGAL")).Block(retVal)
										g.Add(retZero)
									case "string":
										g.If(Id("n").Dot(frag.Name).Op("!=").Lit("")).Block(retVal)
										g.Add(retZero)
									case "bool":
										g.If(Id("n").Dot(frag.Name)).Block(retVal)
										g.Add(retZero)
									default:
										g.Add(retVal)
									}
								}
							})
						}
					}),
				)
			}
		})
		g.Return(Lit(0), Lit(0), Lit(0))
	})

	return f.Save("./decorator/restorer-generated-length.go.txt")
}
