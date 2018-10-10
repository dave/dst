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
```

### Status

This is an experimental package under development, so the API and behaviour is expected to change 
frequently. However I'm now inviting people to try it out and give feedback. 

### Chat?

Feel free to create an [issue](https://github.com/dave/dst/issues) or chat in the 
[#dst](https://gophers.slack.com/messages/CCVL24MTQ) Gophers Slack channel.
