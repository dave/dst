package dst_test

import (
	"go/parser"
	"go/token"

	"go/format"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
)

func Example_Decorations() {
	code := `package main

	func main() {
		var a int
		a++
		print(a)
	}`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "a.go", code, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	df := decorator.Decorate(f, fset)
	df = dstutil.Apply(df, func(c *dstutil.Cursor) bool {
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
	}, nil).(*dst.File)
	f, fset = decorator.Restore(df)
	format.Node(os.Stdout, fset, f)

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
