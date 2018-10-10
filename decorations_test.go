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
			n.Decs.AfterX.Add("/* bar */")
		case *dst.CallExpr:
			n.Decs.AfterLparen.Add("\n")
			n.Decs.AfterArgs.Add("\n")
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

	func main() { 
		var a int    // foo
		var b string // bar
	}`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "a.go", code, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	body := f.Decls[0].(*ast.FuncDecl).Body
	body.List = []ast.Stmt{body.List[1], body.List[0]}

	if err := format.Node(os.Stdout, fset, f); err != nil {
		panic(err)
	}

	//Output:
	//package a
	//
	//func main() {
	//	// foo
	//	var b string
	//	var a int
	//	// bar
	//}
}

func ExampleDstFixed() {
	code := `package a

	func main() { 
		var a int    // foo
		var b string // bar
	}
	`
	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}

	body := f.Decls[0].(*dst.FuncDecl).Body
	body.List = []dst.Stmt{body.List[1], body.List[0]}

	if err := decorator.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package a
	//
	//func main() {
	//	var b string // bar
	//	var a int    // foo
	//}
}
