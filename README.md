# Decorated Syntax Tree

The `dst` package enables manipulation of a Go syntax tree with high fidelity. Decorations (e.g. 
comments and newlines) remain attached to the correct nodes as the tree is modified.

### Where does `go/ast` break?

See [this golang issue](https://github.com/golang/go/issues/20744) for more information.

Consider this example where we want to reverse the order of the two declarations. As you can see the 
comments don't remain attached to the correct nodes:

```go
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
```

Here's the same example using `dst`:

```go
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
```

### Example

This would be prohibitively difficult using `go/ast`:

```go
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
	stmt.Start().Replace(fmt.Sprintf("// foo %d", i))
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
```

### Status

This is an experimental package under development, so the API and behaviour is expected to change 
frequently. However I'm now inviting people to try it out and give feedback. 

### Chat?

Feel free to create an [issue](https://github.com/dave/dst/issues) or chat in the 
[#dst](https://gophers.slack.com/messages/CCVL24MTQ) Gophers Slack channel.
