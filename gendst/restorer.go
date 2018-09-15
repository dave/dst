package main

import (
	"go/types"

	. "github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/loader"
)

func generateRestorer(typeNames []string, nodeInfos map[string]NodeInfo, astPkg *loader.PackageInfo, astTypes map[string]*types.TypeName, dstPkg *loader.PackageInfo, dstTypes map[string]*types.TypeName) error {

	astNode := astPkg.Pkg.Scope().Lookup("Node").Type().Underlying().(*types.Interface)
	dstNode := dstPkg.Pkg.Scope().Lookup("Node").Type().Underlying().(*types.Interface)

	f := NewFile("decorator")
	f.ImportName("github.com/dave/dst", "dst")
	// func (r *Restorer) RestoreNode(n dst.Node) ast.Node {
	// 	switch n := n.(type) {
	// 	case <type>:
	// 		...
	// 	default:
	// 		panic(...)
	// 	}
	// }
	f.Func().Params(Id("r").Op("*").Id("FileRestorer")).Id("RestoreNode").Params(Id("n").Qual("github.com/dave/dst", "Node")).Qual("go/ast", "Node").BlockFunc(func(g *Group) {
		g.If(List(Id("an"), Id("ok")).Op(":=").Id("r").Dot("nodes").Index(Id("n")), Id("ok")).Block(
			Return(Id("an")),
		)
		g.Switch(Id("n").Op(":=").Id("n").Assert(Id("type"))).BlockFunc(func(g *Group) {
			for _, name := range typeNames {
				g.Case(Op("*").Qual("github.com/dave/dst", name)).BlockFunc(func(g *Group) {

					g.Id("r").Dot("applyDecorations").Call(
						Id("n").Dot("Decs"),
						Lit(""),
						Lit(true),
					)

					g.Id("out").Op(":=").Op("&").Qual("go/ast", name).Values()

					nodeInfo := nodeInfos[name]
					posFields := map[string]bool{}
					fragsByName := map[string]FragmentInfo{}
					for _, frag := range nodeInfo.Fragments {
						//fmt.Printf("%#v\n", frag)
						fragsByName[frag.Name] = frag
						if frag.PosField != "" {
							posFields[frag.PosField] = true
						}
					}

					astStruct := astTypes[name].Type().Underlying().(*types.Struct)
					dstStruct := dstTypes[name].Type().Underlying().(*types.Struct)

					dstFields := map[string]*types.Var{}
					for i := 0; i < dstStruct.NumFields(); i++ {
						dstFields[dstStruct.Field(i).Name()] = dstStruct.Field(i)
					}

					for i := 0; i < astStruct.NumFields(); i++ {

						astField := astStruct.Field(i)
						fieldName := astField.Name()

						if name == "File" && fieldName == "Comments" {
							continue
						}

						astFieldTypeName := getTypeName(astField.Type(), astNode)

						dstField, ok := dstFields[astField.Name()]
						var dstFieldTypeName string
						if ok {
							dstFieldTypeName = getTypeName(dstField.Type(), dstNode)
						}
						g.Line()

						frag, hasFrag := fragsByName[fieldName]

						if hasFrag {
							g.Id("r").Dot("applyDecorations").Call(
								Id("n").Dot("Decs"),
								Lit(fieldName),
								Lit(true),
							)
							if frag.PosField != "" {
								if dstFieldTypeName == "bool" {
									// special case - this is handled below
								} else {
									g.Id("out").Dot(frag.PosField).Op("=").Id("r").Dot("cursor")
								}
							}
						}

						switch {
						case astFieldTypeName == "Pos" && dstFieldTypeName == "":
							if fieldName == "From" {
								// Special case "From" / "To" - dst node has "Length"
								g.Id("r").Dot("cursor").Op("+=").Qual("go/token", "Pos").Parens(Id("n").Dot("Length"))
							} else if fieldName == "To" {
								// Special case handled above
							} else {
								if posFields[fieldName] {
									// this is a pos field, so we will render later
								} else {
									g.Id("out").Dot(fieldName).Op("=").Id("r").Dot("cursor")
								}
							}
						case astFieldTypeName == "Node" && dstFieldTypeName == "Node":
							g.If(Id("n").Dot(fieldName).Op("!=").Nil()).BlockFunc(func(g *Group) {
								if hasFrag && frag.PrefixLength > 0 {
									g.Id("r").Dot("cursor").Op("+=").Lit(frag.PrefixLength).Commentf("%s has prefix length %d", fieldName, frag.PrefixLength)
								}
								g.Id("out").Dot(fieldName).Op("=").Id("r").Dot("RestoreNode").Call(Id("n").Dot(fieldName)).Assert(typeLiteral(astField.Type()))
								if hasFrag && frag.SuffixLength > 0 {
									g.Id("r").Dot("cursor").Op("+=").Lit(frag.SuffixLength).Commentf("%s has suffix length %d", fieldName, frag.SuffixLength)
								}
							})
						case astFieldTypeName == "[]Node" && dstFieldTypeName == "[]Node":
							g.For(List(Id("_"), Id("v")).Op(":=").Range().Id("n").Dot(fieldName)).Block(
								Id("out").Dot(fieldName).Op("=").Append(
									Id("out").Dot(fieldName),
									Id("r").Dot("RestoreNode").Call(Id("v")).Assert(typeLiteral(astField.Type().(*types.Slice).Elem())),
								),
							)
						case astFieldTypeName == "Token" && dstFieldTypeName == "Token":
							g.Id("out").Dot(fieldName).Op("=").Id("n").Dot(fieldName)
						case astFieldTypeName == "string" && dstFieldTypeName == "string":
							g.Id("out").Dot(fieldName).Op("=").Id("n").Dot(fieldName)
						case astFieldTypeName == "Pos" && dstFieldTypeName == "bool":
							g.If(Id("n").Dot(fieldName)).Block(
								Id("out").Dot(fieldName).Op("=").Id("r").Dot("cursor"),
							)
						case astFieldTypeName == "ChanDir" && dstFieldTypeName == "ChanDir":
							g.Id("out").Dot(fieldName).Op("=").Qual("go/ast", "ChanDir").Parens(Id("n").Dot(fieldName))
						case astFieldTypeName == "bool" && dstFieldTypeName == "bool":
							g.Id("out").Dot(fieldName).Op("=").Id("n").Dot(fieldName)
						case astFieldTypeName == "Scope" && dstFieldTypeName == "Scope":
							// TODO
							g.Commentf("TODO: %s (Scope)", fieldName)
						case astFieldTypeName == "Object" && dstFieldTypeName == "Object":
							// TODO
							g.Commentf("TODO: %s (Object)", fieldName)
						default:
							g.Commentf("%s %s %s", astField.Name(), astFieldTypeName, dstFieldTypeName)
						}

						if hasFrag {
							if frag.HasLength && frag.Length > 0 {
								g.Id("r").Dot("cursor").Op("+=").Lit(frag.Length).Commentf("%s has fixed length %d", fieldName, frag.Length)
							}
							if frag.LenFieldToken != "" {
								g.Id("r").Dot("cursor").Op("+=").Qual("go/token", "Pos").Parens(Len(Id("n").Dot(frag.LenFieldToken).Dot("String").Call()))
							}
							if frag.LenFieldString != "" {
								g.Id("r").Dot("cursor").Op("+=").Qual("go/token", "Pos").Parens(Len(Id("n").Dot(frag.LenFieldString)))
							}
							g.Id("r").Dot("applyDecorations").Call(
								Id("n").Dot("Decs"),
								Lit(fieldName),
								Lit(false),
							)
						}
					}

					g.Id("r").Dot("applyDecorations").Call(
						Id("n").Dot("Decs"),
						Lit(""),
						Lit(false),
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
