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

func (f *FileRestorer) applyDecorations(decorations dst.Decorations) {
	for _, d := range decorations {

		// for multi-line comments, add a newline for each \n
		if strings.HasPrefix(d, "/*") && strings.Contains(d, "\n") {
			for i, c := range d {
				if c == '\n' {
					f.Lines = append(f.Lines, int(f.cursor)+i)
				}
			}
		}

		// if the decoration is a comment, add it and advance the cursor
		if d != "\n" {
			f.Comments = append(f.Comments, &ast.CommentGroup{List: []*ast.Comment{{Slash: f.cursor, Text: d}}})
			f.cursor += token.Pos(len(d))
		}

		// for newline decorations and also line-comments, add a newline
		if strings.HasPrefix(d, "//") || d == "\n" {
			f.Lines = append(f.Lines, int(f.cursor))
			f.cursor++
		}
	}
}
