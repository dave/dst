// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dstutil_test

import (
	"bytes"

	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
)

var rewriteTests = [...]struct {
	name       string
	orig, want string
	pre, post  dstutil.ApplyFunc
}{
	{name: "nop", orig: "package p\n", want: "package p\n"},

	{name: "replace",
		orig: `package p

var x int
`,
		want: `package p

var t T
`,
		post: func(c *dstutil.Cursor) bool {
			if _, ok := c.Node().(*dst.ValueSpec); ok {
				c.Replace(valspec("t", "T"))
				return false
			}
			return true
		},
	},

	{name: "set doc strings",
		orig: `package p

const z = 0

type T struct{}

var x int
`,
		want: `package p

// a foo is a foo
const z = 0

// a foo is a foo
type T struct{}

// a foo is a foo
var x int
`,
		post: func(c *dstutil.Cursor) bool {
			if gd, ok := c.Node().(*dst.GenDecl); ok {
				gd.Decs.Start.Append("// a foo is a foo")
			}
			return true
		},
	},

	{name: "insert names",
		orig: `package p

const a = 1
`,
		want: `package p

const a, b, c = 1, 2, 3
`,
		pre: func(c *dstutil.Cursor) bool {
			if _, ok := c.Parent().(*dst.ValueSpec); ok {
				switch c.Name() {
				case "Names":
					c.InsertAfter(dst.NewIdent("c"))
					c.InsertAfter(dst.NewIdent("b"))
				case "Values":
					c.InsertAfter(&dst.BasicLit{Kind: token.INT, Value: "3"})
					c.InsertAfter(&dst.BasicLit{Kind: token.INT, Value: "2"})
				}
			}
			return true
		},
	},

	{name: "insert",
		orig: `package p

var (
	x int
	y int
)
`,
		want: `package p

var before1 int
var before2 int
var (
	x int
	y int
)
var after2 int
var after1 int
`,
		pre: func(c *dstutil.Cursor) bool {
			if gd, ok := c.Node().(*dst.GenDecl); ok {
				gd.Decs.Before = dst.NewLine
				c.InsertBefore(vardecl("before1", "int"))
				c.InsertAfter(vardecl("after1", "int"))
				c.InsertAfter(vardecl("after2", "int"))
				c.InsertBefore(vardecl("before2", "int"))
			}
			return true
		},
	},

	{name: "delete",
		orig: `package p

var x int
var y int
var z int
`,
		want: `package p

var y int
var z int
`,
		pre: func(c *dstutil.Cursor) bool {
			n := c.Node()
			if d, ok := n.(*dst.GenDecl); ok && d.Specs[0].(*dst.ValueSpec).Names[0].Name == "x" {
				c.Delete()
			}
			return true
		},
	},

	{name: "insertafter-delete",
		orig: `package p

var x int
var y int
var z int
`,
		want: `package p

var x1 int
var y int
var z int
`,
		pre: func(c *dstutil.Cursor) bool {
			n := c.Node()
			if d, ok := n.(*dst.GenDecl); ok && d.Specs[0].(*dst.ValueSpec).Names[0].Name == "x" {
				c.InsertAfter(vardecl("x1", "int"))
				c.Delete()
			}
			return true
		},
	},

	{name: "delete-insertafter",
		orig: `package p

var x int
var y int
var z int
`,
		want: `package p

var y int
var x1 int
var z int
`,
		pre: func(c *dstutil.Cursor) bool {
			n := c.Node()
			if d, ok := n.(*dst.GenDecl); ok && d.Specs[0].(*dst.ValueSpec).Names[0].Name == "x" {
				c.Delete()
				// The cursor is now effectively atop the 'var y int' node.
				c.InsertAfter(vardecl("x1", "int"))
			}
			return true
		},
	},
}

func valspec(name, typ string) *dst.ValueSpec {
	return &dst.ValueSpec{Names: []*dst.Ident{dst.NewIdent(name)},
		Type: dst.NewIdent(typ),
	}
}

func vardecl(name, typ string) *dst.GenDecl {
	return &dst.GenDecl{
		Tok:   token.VAR,
		Specs: []dst.Spec{valspec(name, typ)},
	}
}

func TestRewrite(t *testing.T) {
	t.Run("*", func(t *testing.T) {
		for _, test := range rewriteTests {
			test := test
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()
				fset := token.NewFileSet()
				f, err := parser.ParseFile(fset, test.name, test.orig, parser.ParseComments)
				if err != nil {
					t.Fatal(err)
				}
				dstFile, err := decorator.DecorateFile(fset, f)
				if err != nil {
					t.Fatal(err)
				}
				dstFile = dstutil.Apply(dstFile, test.pre, test.post).(*dst.File)
				restoredFset, restoredFile, err := decorator.RestoreFile(dstFile)
				if err != nil {
					t.Fatal(err)
				}
				var buf bytes.Buffer
				if err := format.Node(&buf, restoredFset, restoredFile); err != nil {
					t.Fatal(err)
				}
				got := buf.String()
				if got != test.want {
					t.Errorf("got:\n\n%s\nwant:\n\n%s\n", got, test.want)
				}
			})
		}
	})
}

var sink dst.Node

func BenchmarkRewrite(b *testing.B) {
	for _, test := range rewriteTests {
		b.Run(test.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				fset := token.NewFileSet()
				f, err := parser.ParseFile(fset, test.name, test.orig, parser.ParseComments)
				if err != nil {
					b.Fatal(err)
				}
				d, err := decorator.Decorate(fset, f)
				if err != nil {
					b.Fatal(err)
				}
				b.StartTimer()
				sink = dstutil.Apply(d, test.pre, test.post)
			}
		})
	}
}
