package dst_test

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strconv"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"github.com/dave/dst/decorator/resolver/gotypes"
	"github.com/dave/dst/decorator/resolver/guess"
	"github.com/dave/dst/dstutil"
	"golang.org/x/tools/go/packages"
)

func TestDecorations_All(t *testing.T) {
	d := &dst.Decorations{"a", "b"}
	expected := "[a b]"
	found := fmt.Sprint(d.All())
	if expected != found {
		t.Fatalf("expected %s, found %s", expected, found)
	}
}

func TestDecorations_Append(t *testing.T) {
	d := &dst.Decorations{"a", "b"}
	d.Append("c")
	found := fmt.Sprint(*d)
	expected := "[a b c]"
	if expected != found {
		t.Fatalf("expected %s, found %s", expected, found)
	}
}

func TestDecorations_Clear(t *testing.T) {
	d := &dst.Decorations{"a", "b"}
	d.Clear()
	found := fmt.Sprint(*d)
	expected := "[]"
	if expected != found {
		t.Fatalf("expected %s, found %s", expected, found)
	}
}

func TestDecorations_Replace(t *testing.T) {
	d := &dst.Decorations{"a", "b"}
	d.Replace("c")
	found := fmt.Sprint(*d)
	expected := "[c]"
	if expected != found {
		t.Fatalf("expected %s, found %s", expected, found)
	}
}

func TestDecorations_Prepend(t *testing.T) {
	d := &dst.Decorations{"a", "b"}
	d.Prepend("c")
	found := fmt.Sprint(*d)
	expected := "[c a b]"
	if expected != found {
		t.Fatalf("expected %s, found %s", expected, found)
	}
}

func TestSpaceType_String(t *testing.T) {
	if dst.None.String() != "None" {
		t.Fatalf("expected None, found %s", dst.None.String())
	}
	if dst.NewLine.String() != "NewLine" {
		t.Fatalf("expected NewLine, found %s", dst.NewLine.String())
	}
	if dst.EmptyLine.String() != "EmptyLine" {
		t.Fatalf("expected EmptyLine, found %s", dst.EmptyLine.String())
	}
	if dst.SpaceType(99).String() != "" {
		t.Fatalf("expected , found %s", dst.SpaceType(99).String())
	}
}

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
	dir, err := tempDir(map[string]string{
		"go.mod":  "module root",
		"main.go": "package main \n\n func main() {}",
	})
	defer os.RemoveAll(dir)
	if err != nil {
		panic(err)
	}

	// Use the Load convenience function that calls go/packages to load the package. All loaded
	// ast files are decorated to dst.
	pkgs, err := decorator.Load(&packages.Config{Dir: dir, Mode: packages.LoadSyntax}, "root")
	if err != nil {
		panic(err)
	}
	p := pkgs[0]
	f := p.Syntax[0]

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
	if err := r.Print(p.Syntax[0]); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//import "fmt"
	//
	//func main() { fmt.Println("Hello, World!") }
}

func ExampleGoTypesImport() {

	// Create a simple module in a temporary directory
	dir, err := tempDir(map[string]string{
		"go.mod": "module root",
		"main.go": `package main

			import . "fmt" 

			func main() {
				Println("a")
			}`,
	})
	defer os.RemoveAll(dir)
	if err != nil {
		panic(err)
	}

	// Use golang.org/x/tools/go/packages.Load to load the package and parse the syntax
	pkgs, err := packages.Load(&packages.Config{Dir: dir, Mode: packages.LoadSyntax}, "root")
	if err != nil {
		panic(err)
	}
	p := pkgs[0]

	// Create a new decorator to decorate files in the package. This could also be done with the
	// convenience function NewDecoratorFromPackage, but we show the manual method here.
	dec := decorator.NewDecoratorWithImports(p.Fset, p.PkgPath, gotypes.New(p.TypesInfo.Uses))

	f, err := dec.DecorateFile(p.Syntax[0])
	if err != nil {
		panic(err)
	}

	res := decorator.NewRestorerWithImports("root", gopackages.New(dir))
	fr := res.FileRestorer()
	fr.Alias["fmt"] = "" // change the dot-import to a regular import

	if err := fr.Print(f); err != nil {
		panic(err)
	}

	//Output:
	//package main
	//
	//import "fmt"
	//
	//func main() {
	//	fmt.Println("a")
	//}
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

func ExampleDecorationPoints() {
	code := `package main
	
	// main comment
	// is multi line
	func main() {
	
		if true {
	
			// foo
			println( /* foo inline */ "foo")
		} else if false {
			println /* bar inline */ ("bar")
	
			// bar after
	
		} else {
			// empty block
		}
	}`

	f, err := decorator.Parse(code)
	if err != nil {
		panic(err)
	}

	dst.Inspect(f, func(node dst.Node) bool {
		if node == nil {
			return false
		}
		before, after, points := dstutil.Decorations(node)
		var info string
		if before != dst.None {
			info += fmt.Sprintf("- Before: %s\n", before)
		}
		for _, point := range points {
			if len(point.Decs) == 0 {
				continue
			}
			info += fmt.Sprintf("- %s: [", point.Name)
			for i, dec := range point.Decs {
				if i > 0 {
					info += ", "
				}
				info += fmt.Sprintf("%q", dec)
			}
			info += "]\n"
		}
		if after != dst.None {
			info += fmt.Sprintf("- After: %s\n", after)
		}
		if info != "" {
			fmt.Printf("%T\n%s\n", node, info)
		}
		return true
	})

	//Output:
	//*dst.FuncDecl
	//- Before: NewLine
	//- Start: ["// main comment", "// is multi line"]
	//
	//*dst.IfStmt
	//- Before: NewLine
	//- After: NewLine
	//
	//*dst.ExprStmt
	//- Before: NewLine
	//- Start: ["// foo"]
	//- After: NewLine
	//
	//*dst.CallExpr
	//- Lparen: ["/* foo inline */"]
	//
	//*dst.ExprStmt
	//- Before: NewLine
	//- End: ["\n", "\n", "// bar after"]
	//- After: NewLine
	//
	//*dst.CallExpr
	//- Fun: ["/* bar inline */"]
	//
	//*dst.BlockStmt
	//- Lbrace: ["\n", "// empty block"]
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
