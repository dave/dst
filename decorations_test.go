package dst_test

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"github.com/dave/dst/decorator/resolver/guess"
	"golang.org/x/tools/go/packages"
)

func ExampleAlias() {

	code := `package main

		import "fmt"

		func main() {
			fmt.Println("a")
		}`

	dec := decorator.NewDecoratorWithImports(token.NewFileSet(), "main", goast.New())

	f, err := dec.Parse(code)
	if err != nil {
		panic(err)
	}

	res := decorator.NewRestorerWithImports("main", guess.New())

	fr := res.FileRestorer()
	fr.Alias["fmt"] = "fmt1"

	if err := fr.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//import fmt1 "fmt"
	//
	//func main() {
	//	fmt1.Println("a")
	//}

}

func ExampleManualImports() {

	code := `package main

		import "fmt"

		func main() {
			fmt.Println("a")
		}`

	dec := decorator.NewDecoratorWithImports(token.NewFileSet(), "main", goast.New())

	f, err := dec.Parse(code)
	if err != nil {
		panic(err)
	}

	f.Decls[1].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr).Args = []dst.Expr{
		&dst.CallExpr{
			Fun: &dst.Ident{Name: "A", Path: "foo.bar/baz"},
		},
	}

	res := decorator.NewRestorerWithImports("main", guess.New())
	if err := res.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//import (
	//	"fmt"
	//
	//	"foo.bar/baz"
	//)
	//
	//func main() {
	//	fmt.Println(baz.A())
	//}

}

func ExampleImports() {

	// Create a simple module in a temporary directory
	dir, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(dir)
	ioutil.WriteFile(filepath.Join(dir, "go.mod"), []byte("module root"), 0666)
	ioutil.WriteFile(filepath.Join(dir, "main.go"), []byte("package main \n\n func main() {}"), 0666)

	// Use the Load convenience function that calls go/packages to load the package. All loaded
	// ast files are decorated to dst.
	pkgs, err := decorator.Load(&packages.Config{Dir: dir, Mode: packages.LoadSyntax}, "root")
	if err != nil {
		panic(err)
	}
	p := pkgs[0]
	f := p.Files[0]

	// Add a call expression. Note we don't have to use a SelectorExpr - just adding an Ident with
	// the imported package path will do. The restorer will add SelectorExpr where appropriate when
	// converting back to ast. Note the new Path field on *dst.Ident. Set this to the package path
	// of the imported package, and the restorer will automatically add the import to the import
	// block.
	b := f.Decls[0].(*dst.FuncDecl).Body
	b.List = append(b.List, &dst.ExprStmt{
		X: &dst.CallExpr{
			Fun: &dst.Ident{Path: "fmt", Name: "Println"},
			Args: []dst.Expr{
				&dst.BasicLit{Kind: token.STRING, Value: strconv.Quote("Hello, World!")},
			},
		},
	})

	// Create a restorer with the import manager enabled, and print the result. As you can see, the
	// import block is automatically managed, and the Println ident is converted to a SelectorExpr:
	r := decorator.NewRestorerWithImports("root", gopackages.New(dir))
	if err := r.Print(p.Files[0]); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//import "fmt"
	//
	//func main() { fmt.Println("Hello, World!") }
}

func ExampleClone() {
	code := `package main

	var i /* a */ int`

	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}

	cloned := dst.Clone(f.Decls[0]).(*dst.GenDecl)

	cloned.Decs.Before = dst.NewLine
	cloned.Specs[0].(*dst.ValueSpec).Names[0].Name = "j"
	cloned.Specs[0].(*dst.ValueSpec).Names[0].Decs.End.Replace("/* b */")

	f.Decls = append(f.Decls, cloned)

	if err := decorator.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//var i /* a */ int
	//var j /* b */ int
}

func ExampleTypes() {
	code := `package main

	func main() {
		var i int
		i++
		println(i)
	}`

	// Parse the code to AST
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Invoke the type checker using AST as input
	typesInfo := types.Info{
		Defs: make(map[*ast.Ident]types.Object),
		Uses: make(map[*ast.Ident]types.Object),
	}
	conf := &types.Config{}
	if _, err := conf.Check("", fset, []*ast.File{astFile}, &typesInfo); err != nil {
		panic(err)
	}

	// Create a new decorator, which will track the mapping between ast and dst nodes
	dec := decorator.NewDecorator(fset)

	// Decorate the *ast.File to give us a *dst.File
	f, err := dec.DecorateFile(astFile)
	if err != nil {
		panic(err)
	}

	// Find the *dst.Ident for the definition of "i"
	dstDef := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.DeclStmt).Decl.(*dst.GenDecl).Specs[0].(*dst.ValueSpec).Names[0]

	// Find the *ast.Ident using the Ast.Nodes mapping
	astDef := dec.Ast.Nodes[dstDef].(*ast.Ident)

	// Find the types.Object corresponding to "i"
	obj := typesInfo.Defs[astDef]

	// Find all the uses of that object
	var astUses []*ast.Ident
	for id, ob := range typesInfo.Uses {
		if ob != obj {
			continue
		}
		astUses = append(astUses, id)
	}

	// Find each *dst.Ident in the Dst.Nodes mapping
	var dstUses []*dst.Ident
	for _, id := range astUses {
		dstUses = append(dstUses, dec.Dst.Nodes[id].(*dst.Ident))
	}

	// Change the name of the original definition and all uses
	dstDef.Name = "foo"
	for _, id := range dstUses {
		id.Name = "foo"
	}

	// Print the DST
	if err := decorator.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//func main() {
	//	var foo int
	//	foo++
	//	println(foo)
	//}

}

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

	list[0].Decorations().Before = dst.EmptyLine
	list[0].Decorations().End.Append("// the Decorations method allows access to the common")
	list[1].Decorations().End.Append("// decoration properties (Before, Start, End and After)")
	list[2].Decorations().End.Append("// for all nodes.")
	list[2].Decorations().After = dst.EmptyLine

	if err := decorator.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//func main() {
	//
	//	var i int  // the Decorations method allows access to the common
	//	i++        // decoration properties (Before, Start, End and After)
	//	println(i) // for all nodes.
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

	call := f.Decls[0].(*dst.FuncDecl).Body.List[0].(*dst.ExprStmt).X.(*dst.CallExpr)

	call.Decs.Before = dst.EmptyLine
	call.Decs.After = dst.EmptyLine

	for _, v := range call.Args {
		v := v.(*dst.Ident)
		v.Decs.Before = dst.NewLine
		v.Decs.After = dst.NewLine
	}

	if err := decorator.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//func main() {
	//
	//	println(
	//		a,
	//		b,
	//		c,
	//	)
	//
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

	call.Decs.Start.Append("// you can add comments at the start...")
	call.Decs.Fun.Append("/* ...in the middle... */")
	call.Decs.End.Append("// or at the end.")

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
		stmt.Decorations().Before = dst.EmptyLine
		stmt.Decorations().Start.Append(fmt.Sprintf("// foo %d", i))
	}

	call := body.List[2].(*dst.ExprStmt).X.(*dst.CallExpr)
	call.Args = append(call.Args, dst.NewIdent("b"), dst.NewIdent("c"))
	for i, expr := range call.Args {
		expr.Decorations().Before = dst.NewLine
		expr.Decorations().After = dst.NewLine
		expr.Decorations().Start.Append(fmt.Sprintf("/* bar %d */", i))
		expr.Decorations().End.Append(fmt.Sprintf("// baz %d", i))
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
	f, err := parser.ParseFile(fset, "", code, parser.ParseComments)
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
