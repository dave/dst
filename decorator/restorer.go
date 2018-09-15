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

func (f *FileRestorer) applyDecorations(decorations []dst.Decoration, position string, start bool) {
	for _, d := range decorations {
		if d.Position != position || d.Start != start {
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
