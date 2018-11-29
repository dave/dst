package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dave/dst/gendst/data"
	. "github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/loader"
)

// notest

func generateDst(names []string) error {

	f := NewFile("dst")

	f.Comment("notest")
	f.Line()

	for _, name := range names {

		if name == "Package" {
			f.Comment("Decorations is nil for Package nodes.")
		} else {
			f.Comment("Decorations returns the decorations that are common to all nodes (Before, Start, End, After).")
		}
		f.Func().Params(Id("n").Op("*").Id(name)).Id("Decorations").Params().Op("*").Id("NodeDecs").BlockFunc(func(g *Group) {
			if name == "Package" {
				g.Return(Nil())
			} else {
				g.Return(Op("&").Id("n").Dot("Decs").Dot("NodeDecs"))
			}
		})
	}
	return f.Save("./decorations-node-generated.go")
}

func generateDstDecs(names []string) error {

	path := "github.com/dave/dst/gendst/data"
	conf := loader.Config{ParserMode: parser.ParseComments}
	conf.Import(path)
	prog, err := conf.Load()
	if err != nil {
		panic(err)
	}
	var astFile *ast.File
	for _, v := range prog.Package(path).Files {
		_, name := filepath.Split(prog.Fset.File(v.Pos()).Name())
		if name == "positions.go" {
			astFile = v
			break
		}
	}
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
			for _, frag := range data.Info[name] {
				switch frag := frag.(type) {
				case data.Decoration:
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
