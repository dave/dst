package dst_test

import (
	"go/ast"
	"go/parser"
	"go/token"

	"go/format"
	"os"

	"fmt"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

func ExampleSpace() {
	code := `package main

	func main() {
		println(a, b, c)
	}`
	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}

	callExpr := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr)
	callExpr.Decs.Lparen.Add("\n")
	for _, v := range callExpr.Args {
		v.SetSpace(dst.NewLine)
	}
	callExpr.Args[0].End().Add("// you can adjust line-spacing")

	if err := decorator.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//func main() {
	//	println(
	//		a, // you can adjust line-spacing
	//		b,
	//		c,
	//	)
	//}
}

func ExampleComment() {
	code := `package main

	func main() {
		println("Hello World!")
	}`
	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}

	callExpr := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr)

	callExpr.Decs.Start.Add("// you can add comments at the start...")
	callExpr.Decs.Fun.Add("/* ...in the middle... */")
	callExpr.Decs.End.Add("// or at the end.")

	if err := decorator.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//func main() {
	//	// you can add comments at the start...
	//	println /* ...in the middle... */ ("Hello World!") // or at the end.
	//}
}

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

	body := f.Decls[0].(*dst.FuncDecl).Body
	for i, stmt := range body.List {
		stmt.Start().Add(fmt.Sprintf("// foo %d", i))
		stmt.SetSpace(dst.EmptyLine)
	}

	call := body.List[2].(*dst.ExprStmt).X.(*dst.CallExpr)
	call.Args = append(call.Args, dst.NewIdent("b"), dst.NewIdent("c"))
	call.Decs.Lparen.Add("\n")
	for i, expr := range call.Args {
		expr.SetSpace(dst.NewLine)
		expr.Start().Add(fmt.Sprintf("/* bar %d */", i))
		expr.End().Add(fmt.Sprintf("// baz %d", i))
	}

	if err := decorator.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//func main() {
	//	// foo 0
	//	var a int
	//
	//	// foo 1
	//	a++
	//
	//	// foo 2
	//	print(
	//		/* bar 0 */ a, // baz 0
	//		/* bar 1 */ b, // baz 1
	//		/* bar 2 */ c, // baz 2
	//	)
	//
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
