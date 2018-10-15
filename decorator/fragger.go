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

func NewFragger(fset *token.FileSet) *Fragger {
	return &Fragger{
		fset:    fset,
		Indents: map[ast.Node]int{},
	}
}

type Fragger struct {
	cursor    int
	Fragments []Fragment
	Indents   map[ast.Node]int
	fset      *token.FileSet
}

func (f *Fragger) AddDecoration(n ast.Node, name string, pos token.Pos) {
	f.Fragments = append(f.Fragments, &DecorationFragment{Node: n, Name: name, Pos: token.Pos(f.cursor)})
}

func (f *Fragger) AddToken(n ast.Node, t token.Token, pos token.Pos) {
	if pos.IsValid() {
		f.cursor = int(pos)
	}
	f.Fragments = append(f.Fragments, &TokenFragment{Node: n, Token: t, Pos: token.Pos(f.cursor)})
	f.cursor += len(t.String())
}

func (f *Fragger) AddString(n ast.Node, s string, pos token.Pos) {
	if pos.IsValid() {
		f.cursor = int(pos)
	}
	f.Fragments = append(f.Fragments, &StringFragment{Node: n, String: s, Pos: token.Pos(f.cursor)})
	f.cursor += len(s)
}

func (f *Fragger) AddComment(text string, pos token.Pos) {
	// Don't need to worry about the cursor with comments - they are added to the fragment list in
	// the wrong order, then we sort the list based on Pos
	f.Fragments = append(f.Fragments, &CommentFragment{Text: text, Pos: pos})
}

func (f *Fragger) AddNewline(pos token.Pos, empty bool) {
	// Don't need to worry about the cursor with newlines - they are added to the fragment list in
	// the wrong order, then we sort the list based on Pos
	f.Fragments = append(f.Fragments, &NewlineFragment{Pos: pos, Empty: empty})
}

func (f *Fragger) Fragment(node ast.Node) {

	f.ProcessNode(node)

	if f.fset != nil {
		processFile := func(astf *ast.File) {
			// we will avoid adding a newline decoration that is inside a comment
			avoid := map[int]bool{}
			for _, cg := range astf.Comments {
				for _, c := range cg.List {

					// Add the comment to the fragment list.
					f.AddComment(c.Text, c.Slash)

					// Avoid newlines in multi-line comments
					if strings.HasPrefix(c.Text, "/*") {
						startLine := f.fset.Position(c.Pos()).Line
						endLine := f.fset.Position(c.End()).Line

						// multi line comment
						if endLine > startLine {
							for i := startLine; i < endLine; i++ {
								avoid[i+1] = true // we avoid the lines that follow the lines in the comment
							}
						}
					}
				}
			}

			// avoid newlines inside multi-line (back-quoted) strings
			for _, frag := range f.Fragments {
				switch frag := frag.(type) {
				case *StringFragment:
					if !strings.HasPrefix(frag.String, "`") {
						continue
					}

					startLine := f.fset.Position(frag.Pos).Line
					endLine := f.fset.Position(frag.Pos + token.Pos(len(frag.String))).Line

					// multi line string
					if endLine > startLine {
						for i := startLine; i < endLine; i++ {
							avoid[i+1] = true // we avoid the lines that follow the lines in the string
						}
					}
				}
			}

			line := 1
			tokenf := f.fset.File(astf.Pos())
			max := tokenf.Base() + tokenf.Size()
			for i := tokenf.Base(); i < max; i++ {
				pos := f.fset.Position(token.Pos(i))
				if pos.Line != line {

					line = pos.Line

					if avoid[line] {
						continue
					}

					// Peek ahead to the next position in the fset. If we're on another new line, we have
					// an empty line:
					nextLine := line
					if i < max-1 {
						// can't peek forward at the end of the file
						nextLine = f.fset.Position(token.Pos(i + 1)).Line
					}

					if nextLine != line {
						f.AddNewline(token.Pos(i-1), true)
						line = nextLine
						i++
					} else {
						f.AddNewline(token.Pos(i-1), false)
					}

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

	// Search for nodes and comments that start directly after newlines. We note their indent.
	for i, frag := range f.Fragments {
		if i == 0 {
			continue
		}
		if !f.Fragments[i-1].HasNewline() {
			continue
		}
		switch frag := frag.(type) {
		case *DecorationFragment:
			if frag.Name != "Start" {
				continue
			}
			f.Indents[frag.Node] = f.fset.Position(frag.Node.Pos()).Column
		case *CommentFragment:
			frag.Indent = f.fset.Position(frag.Pos).Column
		}
	}
}

func appendDecoration(m map[ast.Node]map[string][]string, n ast.Node, pos, text string) {
	if m[n] == nil {
		m[n] = map[string][]string{}
	}
	m[n][pos] = append(m[n][pos], text)
}

func appendNewLine(m map[ast.Node]map[string][]string, n ast.Node, pos string, empty bool) {
	if m[n] == nil {
		m[n] = map[string][]string{}
	}
	num := 1
	if empty {
		num = 2
	}
	decs := m[n][pos]
	if len(decs) > 0 && strings.HasPrefix(decs[len(decs)-1], "//") {
		num--
	}
	for i := 0; i < num; i++ {
		m[n][pos] = append(m[n][pos], "\n")
	}
}

func (f *Fragger) Link() (space, after map[ast.Node]dst.SpaceType, decorations map[ast.Node]map[string][]string) {

	space = map[ast.Node]dst.SpaceType{}
	after = map[ast.Node]dst.SpaceType{}
	decorations = map[ast.Node]map[string][]string{}

	// Pass 1: associate comment groups with decorations. Sweep up any other comments / new-lines /
	// empty-lines and associate with the same decoration.
	for i, frag := range f.Fragments {
		switch frag := frag.(type) {
		case *CommentFragment:

			if frag.Attached != nil {
				continue
			}

			// Comments (or comment groups) attach to decoration points in this precedence:
			//
			// 1) Before the comment on the same line
			// 2) After the comment on the same line
			// 3) After the comment on subsequent lines (but stopping at empty lines)
			// 4) Before the comment on previous lines (but stopping at empty lines)
			// 5) After the comment on subsequent lines
			// 6) Before the comment on previous lines
			//
			// We always stop at tokens, strings. If we get to the end without finding a decoration point we panic.

			var frags []Fragment // comment / new-line / empty-line
			var dec *DecorationFragment
			var found bool
			var try int
			var onlySearchBackwards = false
			for !found {
				try++
				switch try {
				case 1:
					// Before the comment on the same line (search backwards and stop at any newline)
					frags, dec, found = f.findDecoration(true, true, i, -1, false)
				case 2:
					// Special case for CommClause / CaseClause
					// After the comment on line+2 (search forwards), but ONLY looking for "Start" of
					// CommClause / CaseClause:
					frags1, dec1, found1 := f.findDecoration(false, false, i, 1, true)
					if !found1 {
						continue
					}
					nodeIndent, ok := f.Indents[dec1.Node]
					if !ok {
						// if the node isn't found in Indents, it wasn't at the start of the line. This
						// shouldn't happen for CommClause or CaseClause?
						continue
					}
					if frag.Indent != nodeIndent {
						// The comment is at a different indent to the case clause, so continue but
						// skip all subsequent forward searching steps (we only want to search backwards).
						onlySearchBackwards = true
						continue
					}
					// The comment is at the same indent level as the
					frags = frags1
					dec = dec1
					found = true
				case 3:
					if onlySearchBackwards {
						continue
					}
					// After the comment on the same line
					// After the comment on line+1 (search forwards and stop at any empty line)
					frags, dec, found = f.findDecoration(false, true, i, 1, false)
				case 4:
					// Before the comment on line-1 (search backwards and stop at any empty line)
					frags, dec, found = f.findDecoration(false, true, i, -1, false)
				case 5:
					if onlySearchBackwards {
						continue
					}
					// After the comment on line+2 (search forwards)
					frags, dec, found = f.findDecoration(false, false, i, 1, false)
				case 6:
					// After the comment on line-2 (search backwards)
					frags, dec, found = f.findDecoration(false, false, i, -1, false)
				default:
					panic("no decoration found for " + frag.Text)
				}
			}
			for _, fr := range frags {
				switch fr := fr.(type) {
				case *CommentFragment:
					appendDecoration(decorations, dec.Node, dec.Name, fr.Text)
					fr.Attached = dec
				case *NewlineFragment:
					appendNewLine(decorations, dec.Node, dec.Name, fr.Empty)
					fr.Attached = dec
				}
			}
		}
	}

	// Pass 2: associate any new-lines / empty-lines that have not been added to decorations to node
	// spacing. If they can't be attached as node spacing, attach them as decorations.
	for i, frag := range f.Fragments {
		switch frag := frag.(type) {
		case *NewlineFragment:

			if frag.Attached != nil {
				continue
			}

			// If the newline is directly before / after a node, we can set the Before / After spacing
			// of the node decoration instead of adding the newline as a decoration.
			nodeSpace, _, foundSpace := f.findNode(i, 1)
			nodeAfter, _, foundAfter := f.findNode(i, -1)
			if foundSpace || foundAfter {
				spaceType := dst.NewLine
				if frag.Empty {
					spaceType = dst.EmptyLine
				}
				if foundSpace {
					space[nodeSpace] = spaceType
				}
				if foundAfter {
					after[nodeAfter] = spaceType
				}
				continue
			}

			// If this newline can't be associated with a node, attach it to the next / previous
			// decoration location:
			var dec *DecorationFragment
			var found bool
			var try int
			for !found {
				try++
				switch try {
				case 1:
					// search backwards but stop at any token
					_, dec, found = f.findDecoration(false, false, i, -1, false)
				case 2:
					// search forwards but stop at any token
					_, dec, found = f.findDecoration(false, false, i, 1, false)
				default:
					panic("no decoration found for newline")
				}
			}
			appendNewLine(decorations, dec.Node, dec.Name, frag.Empty)
		}
	}

	return
}

func (f *Fragger) findDecoration(stopAtNewline, stopAtEmptyLine bool, from int, direction int, onlyClause bool) (swept []Fragment, dec *DecorationFragment, found bool) {
	var frags []Fragment
	for i := from; i < len(f.Fragments) && i >= 0; i += direction {
		switch current := f.Fragments[i].(type) {
		case *DecorationFragment:
			if onlyClause {
				switch current.Node.(type) {
				case *ast.CommClause, *ast.CaseClause:
					if current.Name == "Start" {
						return frags, current, true
					}
					return
				default:
					return
				}
			}
			return frags, current, true
		case *NewlineFragment:
			if stopAtNewline {
				return
			}
			if stopAtEmptyLine && current.Empty {
				return
			}
			if current.Attached != nil {
				continue
			}
			if direction == 1 {
				frags = append(frags, current)
			} else {
				frags = append([]Fragment{current}, frags...)
			}
		case *CommentFragment:
			if current.Attached != nil {
				continue
			}
			if direction == 1 {
				frags = append(frags, current)
			} else {
				frags = append([]Fragment{current}, frags...)
			}
		case *TokenFragment, *StringFragment:
			return
		}
	}
	return
}

func (f *Fragger) findNode(from int, direction int) (node ast.Node, dec *DecorationFragment, found bool) {

	var name string
	switch direction {
	case 1:
		name = "Start"
	case -1:
		name = "End"
	}

	for i := from; i < len(f.Fragments) && i >= 0; i += direction {
		switch frag := f.Fragments[i].(type) {
		case *DecorationFragment:
			if frag.Name == name {
				return frag.Node, frag, true
			}
			return
		case *CommentFragment:
			if frag.Attached != nil && frag.Attached.Name == name {
				return frag.Attached.Node, frag.Attached, true
			}
		case *NewlineFragment:
			if frag.Attached != nil && frag.Attached.Name == name {
				return frag.Attached.Node, frag.Attached, true
			}
		case *TokenFragment, *StringFragment:
			return
		}
	}
	return
}

type Fragment interface {
	Position() token.Pos
	HasNewline() bool
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
	Text     string
	Pos      token.Pos
	Attached *DecorationFragment // where did we attach this comment in pass 1?
	Indent   int                 // indent if this comment follows a newline
}

type NewlineFragment struct {
	Pos      token.Pos
	Empty    bool                // true if this newline is an empty line (e.g. follows a "//" comment or "\n")
	Attached *DecorationFragment // where did we attach this comment in pass 1?
}

type DecorationFragment struct {
	Node ast.Node
	Name string
	Pos  token.Pos
}

func (v *TokenFragment) Position() token.Pos      { return v.Pos }
func (v *StringFragment) Position() token.Pos     { return v.Pos }
func (v *CommentFragment) Position() token.Pos    { return v.Pos }
func (v *NewlineFragment) Position() token.Pos    { return v.Pos }
func (v *DecorationFragment) Position() token.Pos { return v.Pos }

func (v *TokenFragment) HasNewline() bool      { return false }
func (v *StringFragment) HasNewline() bool     { return false }
func (v *CommentFragment) HasNewline() bool    { return strings.HasPrefix(v.Text, "//") }
func (v *NewlineFragment) HasNewline() bool    { return true }
func (v *DecorationFragment) HasNewline() bool { return false }

func (f Fragger) debug(fset *token.FileSet, w io.Writer) {
	formatPos := func(s token.Position) string {
		return s.String()[strings.Index(s.String(), ":")+1:]
	}
	nodeType := func(n ast.Node) string {
		return strings.Replace(fmt.Sprintf("%T", n), "*ast.", "", -1)
	}
	for _, v := range f.Fragments {
		switch v := v.(type) {
		case *NewlineFragment:
			if v.Empty {
				fmt.Fprintf(w, "Empty line %s\n", formatPos(fset.Position(v.Pos)))
			} else {
				fmt.Fprintf(w, "New line %s\n", formatPos(fset.Position(v.Pos)))
			}
		case *TokenFragment:
			fmt.Fprintf(w, "%s %q %s\n", nodeType(v.Node), v.Token, formatPos(fset.Position(v.Pos)))
		case *StringFragment:
			fmt.Fprintf(w, "%s %q %s\n", nodeType(v.Node), v.String, formatPos(fset.Position(v.Pos)))
		case *DecorationFragment:
			fmt.Fprintf(w, "%s %s %s\n", nodeType(v.Node), v.Name, formatPos(fset.Position(v.Pos)))
		case *CommentFragment:
			fmt.Fprintf(w, "%q %s\n", v.Text, formatPos(fset.Position(v.Pos)))
		default:
			fmt.Fprintf(w, "%T %s\n", v, formatPos(fset.Position(v.Position())))
		}
	}
}
