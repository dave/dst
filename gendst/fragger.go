package main

import (
	. "github.com/dave/jennifer/jen"
)

func generateFragger(names []string, nodes map[string]NodeInfo) error {
	f := NewFile("decorator")
	f.Func().Params(Id("f").Op("*").Id("Fragger")).Id("ProcessNode").Params(Id("n").Qual("go/ast", "Node")).Block(
		Id("f").Dot("ProcessToken").Call(Id("n"), Lit(""), Id("n").Dot("Pos").Call(), Lit(false)),
		Switch(Id("n").Op(":=").Id("n").Assert(Type())).BlockFunc(func(g *Group) {
			for _, nodeName := range names {
				g.Case(Op("*").Qual("go/ast", nodeName)).BlockFunc(func(g *Group) {

					if nodeName == "FuncDecl" {
						// SPECIAL CASE
						g.Id("f").Dot("funcDeclOverride").Call(Id("n"))
						return
					}

					for _, frag := range nodes[nodeName].Fragments {
						g.Comment(frag.Name)
						switch frag.AstType {
						case "Pos", "Token", "string":
							g.If(Id("n").Dot(frag.AstPositionField).Dot("IsValid").Call()).Block(
								Id("f").Dot("ProcessToken").Call(Id("n"), Lit(frag.Name), Id("n").Dot(frag.AstPositionField), Lit(true)),
							)
						case "Node":
							g.If(Id("n").Dot(frag.Name).Op("!=").Nil()).Block(
								Id("f").Dot("ProcessToken").Call(Id("n"), Lit(frag.Name), Id("n").Dot(frag.Name).Dot("Pos").Call(), Lit(false)),
								Id("f").Dot("ProcessNode").Call(Id("n").Dot(frag.Name)),
								Id("f").Dot("ProcessToken").Call(Id("n"), Lit(frag.Name), Id("n").Dot(frag.Name).Dot("End").Call(), Lit(true)),
							)
						case "[]Node":
							g.If(Id("n").Dot(frag.Name).Op("!=").Nil()).BlockFunc(func(g *Group) {
								//surround := !sub.IsStmt && !sub.IsDecl && sub.Type != "Field"
								surround := false // TODO
								if surround {
									g.Id("f").Dot("ProcessToken").Call(Id("n"), Lit(frag.Name), Id("n").Dot(frag.Name).Index(Lit(0)).Dot("Pos").Call(), Lit(false))
								}
								g.For(List(Id("_"), Id("v")).Op(":=").Range().Id("n").Dot(frag.Name)).Block(
									Id("f").Dot("ProcessNode").Call(Id("v")),
								)
								if surround {
									g.Id("f").Dot("ProcessToken").Call(Id("n"), Lit(frag.Name), Id("n").Dot(frag.Name).Index(Len(Id("n").Dot(frag.Name)).Op("-").Lit(1)).Dot("End").Call(), Lit(true))
								}
							})
						default:
							panic("fragment type " + frag.AstType)
						}
					}
				})
			}
		}),
		Id("f").Dot("ProcessToken").Call(Id("n"), Lit(""), Id("n").Dot("End").Call(), Lit(true)),
	)
	if err := f.Save("./decorator/fragger-generated.go"); err != nil {
		return err
	}
	return nil
}
