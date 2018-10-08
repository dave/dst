package decorator

import (
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"sort"
	"strings"

	"github.com/dave/dst"
)

type Fragger struct {
	cursor    int
	Fragments []Fragment
}

func (f *Fragger) AddDecoration(n ast.Node, name string, pos token.Pos) {
	f.Fragments = append(f.Fragments, DecorationFragment{Node: n, Name: name, Pos: token.Pos(f.cursor)})
}

func (f *Fragger) AddToken(n ast.Node, t token.Token, pos token.Pos) {
	if pos.IsValid() {
		f.cursor = int(pos)
	}
	f.Fragments = append(f.Fragments, TokenFragment{Node: n, Token: t, Pos: token.Pos(f.cursor)})
	f.cursor += len(t.String())
}

func (f *Fragger) AddString(n ast.Node, s string, pos token.Pos) {
	if pos.IsValid() {
		f.cursor = int(pos)
	}
	f.Fragments = append(f.Fragments, StringFragment{Node: n, String: s, Pos: token.Pos(f.cursor)})
	f.cursor += len(s)
}

func (f *Fragger) AddComment(text string, pos token.Pos) {
	// Don't need to worry about the cursor with comments - they are added to the fragment list in
	// the wrong order, then we sort the list based on Pos
	f.Fragments = append(f.Fragments, CommentFragment{Text: text, Pos: pos})
}

func (f *Fragger) AddNewline(pos token.Pos) {
	// Don't need to worry about the cursor with newlines - they are added to the fragment list in
	// the wrong order, then we sort the list based on Pos
	f.Fragments = append(f.Fragments, NewlineFragment{Pos: pos})
}

func (f *Fragger) Fragment(fset *token.FileSet, node ast.Node) {

	f.ProcessNode(node)

	if fset != nil {
		processFile := func(astf *ast.File) {
			// we will avoid adding a newline decoration that is inside a comment
			avoid := map[int]bool{}
			for _, cg := range astf.Comments {
				for _, c := range cg.List {

					// Add the comment to the fragment list.
					f.AddComment(c.Text, c.Slash)

					// Add any newlines.
					if strings.HasPrefix(c.Text, "//") {
						avoid[fset.Position(c.Pos()).Line] = true
					} else if strings.HasPrefix(c.Text, "/*") && fset.Position(c.End()).Line > fset.Position(c.Pos()).Line {
						// multi line comment
						for i := fset.Position(c.Pos()).Line; i < fset.Position(c.End()).Line; i++ {
							avoid[i] = true
						}
					}
				}
			}
			line := 1
			tokenf := fset.File(astf.Pos())
			for i := tokenf.Base(); i < tokenf.Base()+tokenf.Size(); i++ {
				pos := fset.Position(token.Pos(i))
				if pos.Line != line {
					if !avoid[line] {
						f.AddNewline(token.Pos(i - 1))
					}
					line = pos.Line
				}
			}
		}

		switch val := node.(type) {
		case *ast.File:
			processFile(val)
		case *ast.Package:
			for _, file := range val.Files {
				processFile(file)
			}
		}

	}

	sort.SliceStable(f.Fragments, func(i, j int) bool {
		return f.Fragments[i].Position() < f.Fragments[j].Position()
	})

	for i, v := range f.Fragments {
		if i == 0 {
			continue
		}
		switch n := v.(type) {
		case NewlineFragment:
			switch prev := f.Fragments[i-1].(type) {
			case NewlineFragment:
				f.Fragments[i] = NewlineFragment{
					Pos:   n.Pos,
					Empty: true,
				}
			case CommentFragment:
				if strings.HasPrefix(prev.Text, "//") {
					f.Fragments[i] = NewlineFragment{
						Pos:   n.Pos,
						Empty: true,
					}
				}
			}
		}
	}
}

func appendDecoration(m map[ast.Node]map[string][]string, n ast.Node, pos, val string) {
	if m[n] == nil {
		m[n] = map[string][]string{}
	}
	m[n][pos] = append(m[n][pos], val)
}

func (f *Fragger) Link() (before, after map[ast.Node]dst.SpaceType, decorations map[ast.Node]map[string][]string) {

	before = map[ast.Node]dst.SpaceType{}
	after = map[ast.Node]dst.SpaceType{}
	decorations = map[ast.Node]map[string][]string{}

	var lastIndex int
	var lastDecoration DecorationFragment
	for i, frag := range f.Fragments {
		switch frag := frag.(type) {
		case CommentFragment:

			if i < lastIndex {
				// If we've already moved something forward past this index, always move this to the
				// same decoration
				appendDecoration(decorations, lastDecoration.Node, lastDecoration.Name, frag.Text)
				continue
			}

			// Comments (or comment groups) attach to decoration points in this precedence:
			//
			// 1) Before the comment on the same line
			// 2) After the comment on the same line
			// 3) After the comment on line+1
			// 4) Before the comment on line-1
			// 5) After the comment on line+2
			// 6) Before the comment on line-2
			//
			// We always stop at tokens, strings. If we get to the end without finding a decoration point we panic.

			var dec DecorationFragment
			var found bool
			var try int
			var index int
			for !found {
				try++
				switch try {
				case 1:
					// Before the comment on the same line (search backwards and stop at any newline)
					dec, index, found = f.nextDecoration(true, true, i, -1)
				case 2:
					// After the comment on the same line
					// After the comment on line+1 (search forwards and stop at any empty line)
					dec, index, found = f.nextDecoration(false, true, i, 1)
				case 3:
					// Before the comment on line-1 (search backwards and stop at any empty line)
					dec, index, found = f.nextDecoration(false, true, i, -1)
				case 4:
					// After the comment on line+2 (search forwards)
					dec, index, found = f.nextDecoration(false, false, i, 1)
				case 5:
					// After the comment on line-2 (search backwards)
					dec, index, found = f.nextDecoration(false, false, i, -1)
				default:
					panic("no decoration found for " + frag.Text)
				}
			}
			appendDecoration(decorations, dec.Node, dec.Name, frag.Text)
			lastIndex = index
			lastDecoration = dec

		case NewlineFragment:

			if i < lastIndex {
				// If we've already moved something forward past this index, always move this to the
				// same decoration
				appendDecoration(decorations, lastDecoration.Node, lastDecoration.Name, "\n")
				continue
			}

			// If the newline is directly before / after a node, we can set the Before / After spacing
			// of the node decoration instead of adding the newline as a decoration.
			// TODO...

			var dec DecorationFragment
			var found bool
			var try int
			var index int
			for !found {
				try++
				switch try {
				case 1:
					// search backwards but stop at any token
					dec, index, found = f.nextDecoration(false, false, i, -1)
				case 2:
					// search forwards but stop at any token
					dec, index, found = f.nextDecoration(false, false, i, 1)
				default:
					panic("no decoration found for newline")
				}
			}
			appendDecoration(decorations, dec.Node, dec.Name, "\n")
			lastIndex = index
			lastDecoration = dec
		}
	}
	return
}

func (f *Fragger) nextDecoration(stopAtNewline, stopAtEmptyLine bool, from int, direction int) (frag DecorationFragment, index int, found bool) {
	for i := from; i < len(f.Fragments) && i >= 0; i += direction {
		switch f := f.Fragments[i].(type) {
		case DecorationFragment:
			return f, i, true
		case NewlineFragment:
			if stopAtNewline {
				return
			}
			if stopAtEmptyLine && f.Empty {
				return
			}
		case TokenFragment, StringFragment:
			return
		}
	}
	return
}

//func (f *Fragger) findNodeEnds(from int, direction int) (frag DecorationFragment, index int, found bool) {

//}

type Fragment interface {
	Position() token.Pos
}

type TokenFragment struct {
	Node  ast.Node
	Token token.Token
	Pos   token.Pos
}

type StringFragment struct {
	Node   ast.Node
	String string
	Pos    token.Pos
}

type CommentFragment struct {
	Text string
	Pos  token.Pos
}

type NewlineFragment struct {
	Pos   token.Pos
	Empty bool // true if this newline is an empty line (e.g. follows a "//" comment or "\n")
}

type DecorationFragment struct {
	Node ast.Node
	Name string
	Pos  token.Pos
}

func (v TokenFragment) Position() token.Pos      { return v.Pos }
func (v StringFragment) Position() token.Pos     { return v.Pos }
func (v CommentFragment) Position() token.Pos    { return v.Pos }
func (v NewlineFragment) Position() token.Pos    { return v.Pos }
func (v DecorationFragment) Position() token.Pos { return v.Pos }

func (f Fragger) debug(fset *token.FileSet, w io.Writer) {
	formatPos := func(s token.Position) string {
		return s.String()[strings.Index(s.String(), ":")+1:]
	}
	nodeType := func(n ast.Node) string {
		return strings.Replace(fmt.Sprintf("%T", n), "*ast.", "", -1)
	}
	for _, v := range f.Fragments {
		switch v := v.(type) {
		case NewlineFragment:
			if v.Empty {
				fmt.Fprintf(w, "Empty line %s\n", formatPos(fset.Position(v.Pos)))
			} else {
				fmt.Fprintf(w, "New line %s\n", formatPos(fset.Position(v.Pos)))
			}
		case TokenFragment:
			fmt.Fprintf(w, "%s %q %s\n", nodeType(v.Node), v.Token, formatPos(fset.Position(v.Pos)))
		case StringFragment:
			fmt.Fprintf(w, "%s %q %s\n", nodeType(v.Node), v.String, formatPos(fset.Position(v.Pos)))
		case DecorationFragment:
			fmt.Fprintf(w, "%s %s %s\n", nodeType(v.Node), v.Name, formatPos(fset.Position(v.Pos)))
		case CommentFragment:
			fmt.Fprintf(w, "%q %s\n", v.Text, formatPos(fset.Position(v.Pos)))
		default:
			fmt.Fprintf(w, "%T %s\n", v, formatPos(fset.Position(v.Position())))
		}
	}
}
