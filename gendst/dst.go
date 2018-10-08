package main

import (
	"go/parser"
	"regexp"

	"bytes"
	"go/format"

	"strings"

	"github.com/dave/dst/gendst/fragment"
	. "github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/loader"
)

func generateDst(names []string) error {

	names = append(names, "DecorationStmt", "DecorationDecl")

	f := NewFile("dst")
	for _, name := range names {
		f.Func().Params(Id("v").Op("*").Id(name)).Id("isNode").Params().Block()
	}
	return f.Save("./generated.go")
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
	if len(parts) > 0 {
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
			g.Id("Before").Id("SpaceType")
			for _, frag := range fragment.Info[name] {
				switch frag := frag.(type) {
				case fragment.Decoration:
					g.Id(frag.Name).Id("Decorations")
				}
			}
			g.Id("After").Id("SpaceType")
		})
	}
	return f.Save("./generated-decs.go")
}
