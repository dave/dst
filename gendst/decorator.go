package main

import (
	"fmt"
	"go/types"

	. "github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/loader"
)

func generateDecorator(typeNames []string, astPkg *loader.PackageInfo, astTypes map[string]*types.TypeName, dstPkg *loader.PackageInfo, dstTypes map[string]*types.TypeName) error {

	astNode := astPkg.Pkg.Scope().Lookup("Node").Type().Underlying().(*types.Interface)
	dstNode := dstPkg.Pkg.Scope().Lookup("Node").Type().Underlying().(*types.Interface)

	f := NewFile("decorator")
	f.ImportName("github.com/dave/dst", "dst")

	/*
		func (d *Decorator) NodeToDst(n ast.Node) dst.Node {
			if dn, ok := d.nodes[n]; ok {
				return dn
			}
		}
	*/
	f.Func().Params(Id("d").Op("*").Id("Decorator")).Id("NodeToDst").Params(
		Id("n").Qual("go/ast", "Node"),
	).Qual("github.com/dave/dst", "Node").BlockFunc(func(g *Group) {
		g.If(List(Id("dn"), Id("ok")).Op(":=").Id("d").Dot("nodes").Index(Id("n")), Id("ok")).Block(
			Return(Id("dn")),
		)
		g.Switch(Id("n").Op(":=").Id("n").Assert(Id("type"))).BlockFunc(func(g *Group) {
			for _, name := range typeNames {
				astStruct := astTypes[name].Type().Underlying().(*types.Struct)
				dstStruct := dstTypes[name].Type().Underlying().(*types.Struct)
				dstFields := map[string]*types.Var{}
				for i := 0; i < dstStruct.NumFields(); i++ {
					dstFields[dstStruct.Field(i).Name()] = dstStruct.Field(i)
				}
				g.Case(Op("*").Qual("go/ast", name)).BlockFunc(func(g *Group) {
					g.Id("out").Op(":=").Op("&").Qual("github.com/dave/dst", name).Values()
					for i := 0; i < astStruct.NumFields(); i++ {
						astField := astStruct.Field(i)
						dstField, ok := dstFields[astField.Name()]
						if ok {
							astFieldTypeName := getTypeName(astField.Type(), astNode)
							dstFieldTypeName := getTypeName(dstField.Type(), dstNode)

							switch {
							case astFieldTypeName == "Node" && dstFieldTypeName == "Node":
								// both nodes, just recurse
								g.If(Id("n").Dot(dstField.Name()).Op("!=").Nil()).Block(
									Id("out").Dot(dstField.Name()).Op("=").Id("d").Dot("NodeToDst").Call(
										Id("n").Dot(astField.Name()),
									).Assert(typeLiteral(dstField.Type())),
								)
							case astFieldTypeName == "[]Node" && dstFieldTypeName == "[]Node":
								g.For(List(Id("_"), Id("v")).Op(":=").Range().Id("n").Dot(astField.Name())).Block(
									Id("out").Dot(dstField.Name()).Op("=").Append(
										Id("out").Dot(dstField.Name()),
										Id("d").Dot("NodeToDst").Call(Id("v")).Assert(typeLiteral(dstField.Type().(*types.Slice).Elem())),
									),
								)
							case astFieldTypeName == "string" && dstFieldTypeName == "string":
								g.Id("out").Dot(dstField.Name()).Op("=").Id("n").Dot(astField.Name())
							case astFieldTypeName == "bool" && dstFieldTypeName == "bool":
								g.Id("out").Dot(dstField.Name()).Op("=").Id("n").Dot(astField.Name())
							case astFieldTypeName == "Token" && dstFieldTypeName == "Token":
								g.Id("out").Dot(dstField.Name()).Op("=").Id("n").Dot(astField.Name())
							case astFieldTypeName == "ChanDir" && dstFieldTypeName == "ChanDir":
								g.Id("out").Dot(dstField.Name()).Op("=").Qual("github.com/dave/dst", "ChanDir").Parens(Id("n").Dot(astField.Name()))
							case astFieldTypeName == "Pos" && dstFieldTypeName == "bool":
								g.If(Id("n").Dot(astField.Name()).Op("!=").Qual("go/token", "NoPos")).Block(
									Id("out").Dot(dstField.Name()).Op("=").True(),
								)
							case astFieldTypeName == "Scope" && dstFieldTypeName == "Scope":
								// TODO
								g.Commentf("TODO: %s (Scope)", astField.Name())
							case astFieldTypeName == "Object" && dstFieldTypeName == "Object":
								// TODO
								g.Commentf("TODO: %s (Object)", astField.Name())
							default:
								panic(fmt.Sprintf("%s: ast: %s, dst: %s", astField.Name(), astFieldTypeName, dstFieldTypeName))
							}
						}
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

	if err := f.Save("./decorator/decorator-generated.go"); err != nil {
		return nil
	}
	return nil
}

func typeLiteral(t types.Type) *Statement {
	if p, ok := t.(*types.Pointer); ok {
		return Op("*").Qual(
			p.Elem().(*types.Named).Obj().Pkg().Path(),
			p.Elem().(*types.Named).Obj().Name(),
		)
	} else {
		return Qual(
			t.(*types.Named).Obj().Pkg().Path(),
			t.(*types.Named).Obj().Name(),
		)
	}
}
