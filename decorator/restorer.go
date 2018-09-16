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
	fsetFile := r.Fset.AddFile(fname, r.Fset.Base(), int(fr.cursor-fr.base))
	for _, cg := range fr.Comments {
		astFile.Comments = append(astFile.Comments, cg)
	}
	fsetFile.SetLines(fr.Lines)

	return astFile
}

func (f *FileRestorer) applyDecorations(decorations []dst.Decoration, position string, end bool) {
	for _, d := range decorations {
		if d.Position != position || d.End != end {
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

func (f *FileRestorer) customLength(node dst.Node, fragment string) (suffix, length, prefix int) {

	// value == -1 => calculate from field

	switch node := node.(type) {
	case *dst.ArrayType:
		switch fragment {
		case "Len":
			return 0, -1, 1 // Len has "]" suffix even when Len == nil
		}
	case *dst.SliceExpr:
		switch fragment {
		case "Max":
			if node.Slice3 {
				// If Slice3, we have two colons even with Max == nil
				return 1, -1, 0
			} else if node.Max != nil {
				return 1, -1, 0
			} else {
				return 0, -1, 0
			}
		}
	case *dst.ChanType:
		switch fragment {
		case "Arrow":
			if node.Dir == 0 {
				return 0, 0, 0
			} else {
				return 0, 2, 0
			}
		}
	case *dst.EmptyStmt:
		// TODO: Is this needed?
		if node.Implicit {
			return 0, 0, 0
		} else {
			return 0, 1, 0
		}
	case *dst.TypeAssertExpr:
		switch fragment {
		case "Type":
			if node.Type == nil {
				return 0, len(".(type)"), 0
			} else {
				return len(".("), -1, len(")")
			}
		}
	case *dst.CommClause:
		switch fragment {
		case "Comm":
			if node.Comm == nil {
				return 0, len("default"), 0
			} else {
				return len("case"), -1, 0
			}
		}
	}
	return -1, -1, -1
}

func (r *FileRestorer) funcDeclOverride(n *dst.FuncDecl) *ast.FuncDecl {
	r.applyDecorations(n.Decs, "", false)
	out := &ast.FuncDecl{}
	{
		r.applyDecorations(n.Decs, "Doc", false)
		prefix, length, suffix := getLength(n, "Doc")
		r.cursor += token.Pos(prefix)
		if n.Doc != nil {
			out.Doc = r.RestoreNode(n.Doc).(*ast.CommentGroup)
		}
		r.cursor += token.Pos(length)
		r.cursor += token.Pos(suffix)
		r.applyDecorations(n.Decs, "Doc", true)
	}
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
