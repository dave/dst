package main

import (
	"bytes"
	"go/format"
	"go/parser"
	"regexp"
	"strings"

	"github.com/dave/dst/gendst/fragment"
	. "github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/loader"
)

func generateDst(names []string) error {

	f := NewFile("dst")
	for _, name := range names {

		if name == "Package" {
			f.Comment("Decorations is nil for Package nodes.")
		} else {
			f.Comment("Decorations returns the decorations that are common to all nodes (Space, Start, End, After).")
		}
		f.Func().Params(Id("n").Op("*").Id(name)).Id("Decorations").Params().Op("*").Id("NodeDecs").BlockFunc(func(g *Group) {
			if name == "Package" {
				g.Return(Nil())
			} else {
				g.Return(Op("&").Id("n").Dot("Decs").Dot("NodeDecs"))
			}
		})

		/*
			for _, frag := range fragment.Info[name] {
				switch frag := frag.(type) {
				case fragment.Decoration:
					if frag.Name == "Start" {
						f.Func().Params(Id("v").Op("*").Id(name)).Id("Start").Params().Op("*").Id("Decorations").Block(
							Return(Op("&").Id("v").Dot("Decs").Dot("Start")),
						)
					}
					if frag.Name == "End" {
						f.Func().Params(Id("v").Op("*").Id(name)).Id("End").Params().Op("*").Id("Decorations").Block(
							Return(Op("&").Id("v").Dot("Decs").Dot("End")),
						)
					}
				}
			}

			if name != "Package" {
				f.Func().Params(Id("v").Op("*").Id(name)).Id("Space").Params().Id("SpaceType").Block(
					Return(Id("v").Dot("Decs").Dot("Space")),
				)
				f.Func().Params(Id("v").Op("*").Id(name)).Id("SetSpace").Params(Id("s").Id("SpaceType")).Block(
					Id("v").Dot("Decs").Dot("Space").Op("=").Id("s"),
				)
				f.Func().Params(Id("v").Op("*").Id(name)).Id("After").Params().Id("SpaceType").Block(
					Return(Id("v").Dot("Decs").Dot("After")),
				)
				f.Func().Params(Id("v").Op("*").Id(name)).Id("SetAfter").Params(Id("s").Id("SpaceType")).Block(
					Id("v").Dot("Decs").Dot("After").Op("=").Id("s"),
				)
			}
		*/
	}
	return f.Save("./decorations-node-generated.go")
}

func generateDstDecs(names []string) error {

	path := "github.com/dave/dst/gendst/postests"
	conf := loader.Config{ParserMode: parser.ParseComments}
	conf.Import(path)
	prog, err := conf.Load()
	if err != nil {
		panic(err)
	}
	astFile := prog.Package(path).Files[0]
	buf := &bytes.Buffer{}
	if err := format.Node(buf, prog.Fset, astFile); err != nil {
		panic(err)
	}
	source := buf.String()
	reg := regexp.MustCompile(`// ([a-zA-Z]+)`)
	type part struct {
		name       string
		start, end int
	}
	var parts []part
	for _, cg := range astFile.Comments {
		for _, c := range cg.List {
			if strings.HasPrefix(c.Text, "// --") {
				if len(parts) > 0 && parts[len(parts)-1].end == -1 {
					parts[len(parts)-1].end = int(c.Pos() - 1)
				}
				continue
			}
			if matches := reg.FindStringSubmatch(c.Text); matches != nil {
				name := matches[1]
				pos := c.End()
				prev := c.Pos() - 1
				if len(parts) > 0 && parts[len(parts)-1].end == -1 {
					parts[len(parts)-1].end = int(prev)
				}
				parts = append(parts, part{name, int(pos), -1})
			}
		}
	}
	if len(parts) > 0 && parts[len(parts)-1].end == -1 {
		parts[len(parts)-1].end = int(astFile.End())
	}

	f := NewFile("dst")
	for _, name := range names {
		// type <name>Decorations struct {
		// 	//...
		// }
		f.Line()
		f.Commentf("%sDecorations holds decorations for %s:", name, name)
		f.Comment("")
		for _, part := range parts {
			if part.name != name {
				continue
			}
			text := source[part.start:part.end]
			indented := text[0] == '\t'
			text = strings.TrimSpace(text)
			var indent string
			if !indented || name == "LabeledStmt" { // LabeledStmt special case because comment is in wrong position
				indent = "\t"
			}
			text = "// \t" + strings.Replace(text, "\n", "\n// "+indent, -1)
			f.Comment(text)
			f.Comment("")
		}
		f.Type().Id(name + "Decorations").StructFunc(func(g *Group) {
			g.Id("NodeDecs")
			for _, frag := range fragment.Info[name] {
				switch frag := frag.(type) {
				case fragment.Decoration:
					if frag.Name == "Start" || frag.Name == "End" {
						continue
					}
					g.Id(frag.Name).Id("Decorations")
				}
			}
		})
	}
	return f.Save("./decorations-types-generated.go")
}
