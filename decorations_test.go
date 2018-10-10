package dst_test

import (
	"go/ast"
	"go/parser"
	"go/token"

	"go/format"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
)

func ExampleDecorations() {
	code := `package main

	func main() {
		var a int
		a++
		print(a)
	}`
	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}
	apply := func(c *dstutil.Cursor) bool {
		switch n := c.Node().(type) {
		case *dst.DeclStmt:
			n.Decs.End.Replace("// foo")
		case *dst.IncDecStmt:
			n.Decs.X.Add("/* bar */")
		case *dst.CallExpr:
			n.Decs.Lparen.Add("\n")
			n.Decs.Args.Add("\n")
		}
		return true
	}
	f = dstutil.Apply(f, apply, nil).(*dst.File)
	if err := decorator.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//func main() {
	//	var a int // foo
	//	a /* bar */ ++
	//	print(
	//		a,
	//	)
	//}
}

func ExampleAstBroken() {
	code := `package a

	var a int    // foo
	var b string // bar
	`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "a.go", code, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	f.Decls = []ast.Decl{f.Decls[1], f.Decls[0]}

	if err := format.Node(os.Stdout, fset, f); err != nil {
		panic(err)
	}

	//Output:
	//package a
	//
	//// foo
	//var b string
	//var a int
	//
	//// bar
}

func ExampleDstFixed() {
	code := `package a

	var a int    // foo
	var b string // bar
	`
	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}

	f.Decls = []dst.Decl{f.Decls[1], f.Decls[0]}

	if err := decorator.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package a
	//
	//var b string // bar
	//var a int    // foo
}
