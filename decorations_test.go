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

func ExampleDecorated() {
	code := `package main

	func main() {
		var i int
		i++
		println(i)
	}`
	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}

	list := f.Decls[0].(*dst.FuncDecl).Body.List

	list[0].SetSpace(dst.EmptyLine)
	list[0].End().Add("// the Decorated interface allows access to the common")
	list[1].End().Add("// decoration properties (Space, Start, End and After)")
	list[2].End().Add("// for all Expr, Stmt and Decl nodes.")
	list[2].SetAfter(dst.EmptyLine)

	if err := decorator.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//func main() {
	//
	//	var i int  // the Decorated interface allows access to the common
	//	i++        // decoration properties (Space, Start, End and After)
	//	println(i) // for all Expr, Stmt and Decl nodes.
	//
	//}
}

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
	for i, v := range callExpr.Args {
		switch v := v.(type) {
		case *dst.Ident:
			if i == 0 {
				v.Decs.End.Add("// you can adjust line-spacing")
			}
			v.Decs.Space = dst.NewLine
			v.Decs.After = dst.NewLine
		}
	}

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

	call := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr)

	call.Decs.Start.Add("// you can add comments at the start...")
	call.Decs.Fun.Add("/* ...in the middle... */")
	call.Decs.End.Add("// or at the end.")

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
		stmt.SetSpace(dst.EmptyLine)
		stmt.Start().Add(fmt.Sprintf("// foo %d", i))
	}

	call := body.List[2].(*dst.ExprStmt).X.(*dst.CallExpr)
	call.Args = append(call.Args, dst.NewIdent("b"), dst.NewIdent("c"))
	for i, expr := range call.Args {
		expr.SetSpace(dst.NewLine)
		expr.SetAfter(dst.NewLine)
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
	//
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
	//}
}

func ExampleAstBroken() {
	code := `package a

	func main(){
		var a int    // foo
		var b string // bar
	}
	`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "a.go", code, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	list := f.Decls[0].(*ast.FuncDecl).Body.List
	list[0], list[1] = list[1], list[0]

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

	func main(){
		var a int    // foo
		var b string // bar
	}
	`
	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}

	list := f.Decls[0].(*dst.FuncDecl).Body.List
	list[0], list[1] = list[1], list[0]

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
