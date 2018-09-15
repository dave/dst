package main

import (
	. "github.com/dave/jennifer/jen"
)

func generateProcessor(nodeTypes []string, nodeInfos map[string]NodeInfo) error {
	f := NewFile("decorator")
	f.Func().Params(Id("f").Op("*").Id("Fragger")).Id("ProcessNode").Params(Id("n").Qual("go/ast", "Node")).Block(
		Id("f").Dot("ProcessToken").Call(Id("n"), Lit(""), Lit(true), Lit(0), Id("n").Dot("Pos").Call()),
		Switch(Id("n").Op(":=").Id("n").Assert(Type())).BlockFunc(func(g *Group) {
			for _, nodeType := range nodeTypes {
				g.Case(Op("*").Qual("go/ast", nodeType)).BlockFunc(func(g *Group) {
					for _, frag := range nodeInfos[nodeType].Fragments {
						g.Comment(frag.Name)
						switch {
						case frag.Type == "Pos", frag.Type == "Token", frag.Type == "String":
							var length *Statement
							if frag.HasLength {
								length = Lit(frag.Length)
							} else if frag.LenFieldString != "" {
								length = Len(Id("n").Dot(frag.LenFieldString))
							} else if frag.LenFieldToken != "" {
								length = Len(Id("n").Dot(frag.LenFieldToken).Dot("String").Call())
							}
							g.If(Id("n").Dot(frag.PosField).Dot("IsValid").Call()).Block(
								Id("f").Dot("ProcessToken").Call(Id("n"), Lit(frag.Name), Lit(false), length, Id("n").Dot(frag.PosField)),
							)
						case frag.IsNode:
							g.If(Id("n").Dot(frag.Name).Op("!=").Nil()).BlockFunc(func(g *Group) {
								if frag.Slice {
									//surround := !sub.IsStmt && !sub.IsDecl && sub.Type != "Field"
									surround := false // TODO
									if surround {
										g.Id("f").Dot("ProcessToken").Call(Id("n"), Lit(frag.Name), Lit(true), Lit(0), Id("n").Dot(frag.Name).Index(Lit(0)).Dot("Pos").Call())
									}
									g.For(List(Id("_"), Id("v")).Op(":=").Range().Id("n").Dot(frag.Name)).Block(
										Id("f").Dot("ProcessNode").Call(Id("v")),
									)
									if surround {
										g.Id("f").Dot("ProcessToken").Call(Id("n"), Lit(frag.Name), Lit(false), Lit(0), Id("n").Dot(frag.Name).Index(Len(Id("n").Dot(frag.Name)).Op("-").Lit(1)).Dot("End").Call())
									}
								} else {
									g.Id("f").Dot("ProcessToken").Call(Id("n"), Lit(frag.Name), Lit(true), Lit(0), Id("n").Dot(frag.Name).Dot("Pos").Call())
									g.Id("f").Dot("ProcessNode").Call(Id("n").Dot(frag.Name))
									g.Id("f").Dot("ProcessToken").Call(Id("n"), Lit(frag.Name), Lit(false), Lit(0), Id("n").Dot(frag.Name).Dot("End").Call())
								}

							})
						}
					}
				})
			}
		}),
		Id("f").Dot("ProcessToken").Call(Id("n"), Lit(""), Lit(false), Lit(0), Id("n").Dot("End").Call()),
	)
	if err := f.Save("./decorator/fragger-generated.go"); err != nil {
		return err
	}
	return nil
}
