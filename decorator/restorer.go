package decorator

import (
	"go/ast"
	"go/token"

	"strings"

	"github.com/dave/dst"
)

type Restorer struct {
	Fset *token.FileSet
}

type FileRestorer struct {
	*Restorer
	Lines    []int
	Comments []*ast.CommentGroup
	base     token.Pos
	cursor   token.Pos
	nodes    map[dst.Node]ast.Node
}

func (r *Restorer) Restore(fname string, dstFile *dst.File) *ast.File {
	if r.Fset == nil {
		r.Fset = token.NewFileSet()
	}
	fr := &FileRestorer{
		Restorer: r,
		base:     token.Pos(r.Fset.Base()),
		cursor:   token.Pos(r.Fset.Base()),
		nodes:    map[dst.Node]ast.Node{},
		Lines:    []int{0},
	}
	astFile := fr.RestoreNode(dstFile).(*ast.File)

	astFileEndPos := int(fr.cursor - fr.base)
	// Check that none of the comments or newlines extend past the file end position. If so, increment.
	for _, cg := range fr.Comments {
		if int(cg.End()) >= astFileEndPos {
			astFileEndPos = int(cg.End()) + 1
		}
	}
	for _, l := range fr.Lines {
		if l >= astFileEndPos {
			astFileEndPos = l + 1
		}
	}

	fsetFile := r.Fset.AddFile(fname, r.Fset.Base(), astFileEndPos)
	for _, cg := range fr.Comments {
		astFile.Comments = append(astFile.Comments, cg)
	}
	success := fsetFile.SetLines(fr.Lines)
	if !success {
		panic("SetLines failed")
	}

	return astFile
}

func (f *FileRestorer) applyDecorations(decorations []dst.Decoration, position string, end bool) {
	for _, d := range decorations {
		if d.Position != position {
			continue
		}

		// for multi-line comments, add a newline for each \n
		if strings.HasPrefix(d.Text, "/*") && strings.Contains(d.Text, "\n") {
			for i, c := range d.Text {
				if c == '\n' {
					f.Lines = append(f.Lines, int(f.cursor)+i)
				}
			}
		}

		// if the decoration is a comment, add it and advance the cursor
		if d.Text != "\n" {
			f.Comments = append(f.Comments, &ast.CommentGroup{List: []*ast.Comment{{Slash: f.cursor, Text: d.Text}}})
			f.cursor += token.Pos(len(d.Text))
		}

		// for newline decorations and also line-comments, add a newline
		if strings.HasPrefix(d.Text, "//") || d.Text == "\n" {
			f.Lines = append(f.Lines, int(f.cursor))
			f.cursor++
		}
	}
}

func (r *FileRestorer) funcDeclOverride(n *dst.FuncDecl) *ast.FuncDecl {
	r.applyDecorations(n.Decs, "", false)
	out := &ast.FuncDecl{}
	out.Type = &ast.FuncType{}
	r.nodes[n.Type] = out.Type
	{
		r.applyDecorations(n.Decs, "Func", false)
		prefix, length, suffix := getLength(n.Type, "Func")
		if n.Type.Func {
			out.Type.Func = r.cursor
		}
		r.cursor += token.Pos(prefix)
		r.cursor += token.Pos(length)
		r.cursor += token.Pos(suffix)
		r.applyDecorations(n.Decs, "Func", true)
	}
	{
		r.applyDecorations(n.Decs, "Recv", false)
		prefix, length, suffix := getLength(n, "Recv")
		r.cursor += token.Pos(prefix)
		if n.Recv != nil {
			out.Recv = r.RestoreNode(n.Recv).(*ast.FieldList)
		}
		r.cursor += token.Pos(length)
		r.cursor += token.Pos(suffix)
		r.applyDecorations(n.Decs, "Recv", true)
	}
	{
		r.applyDecorations(n.Decs, "Name", false)
		prefix, length, suffix := getLength(n, "Name")
		r.cursor += token.Pos(prefix)
		if n.Name != nil {
			out.Name = r.RestoreNode(n.Name).(*ast.Ident)
		}
		r.cursor += token.Pos(length)
		r.cursor += token.Pos(suffix)
		r.applyDecorations(n.Decs, "Name", true)
	}
	{
		r.applyDecorations(n.Decs, "Params", false)
		prefix, length, suffix := getLength(n.Type, "Params")
		r.cursor += token.Pos(prefix)
		if n.Type.Params != nil {
			out.Type.Params = r.RestoreNode(n.Type.Params).(*ast.FieldList)
		}
		r.cursor += token.Pos(length)
		r.cursor += token.Pos(suffix)
		r.applyDecorations(n.Decs, "Params", true)
	}
	{
		r.applyDecorations(n.Decs, "Results", false)
		prefix, length, suffix := getLength(n.Type, "Results")
		r.cursor += token.Pos(prefix)
		if n.Type.Results != nil {
			out.Type.Results = r.RestoreNode(n.Type.Results).(*ast.FieldList)
		}
		r.cursor += token.Pos(length)
		r.cursor += token.Pos(suffix)
		r.applyDecorations(n.Decs, "Results", true)
	}
	{
		r.applyDecorations(n.Decs, "Body", false)
		prefix, length, suffix := getLength(n, "Body")
		r.cursor += token.Pos(prefix)
		if n.Body != nil {
			out.Body = r.RestoreNode(n.Body).(*ast.BlockStmt)
		}
		r.cursor += token.Pos(length)
		r.cursor += token.Pos(suffix)
		r.applyDecorations(n.Decs, "Body", true)
	}
	r.applyDecorations(n.Decs, "", true)
	r.nodes[n] = out
	return out
}
