package decorator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/dave/dst"
)

func Parse(src interface{}) (*dst.File, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return Decorate(fset, f).(*dst.File), nil
}

func ParseFile(fset *token.FileSet, filename string, src interface{}, mode parser.Mode) (*dst.File, error) {
	f, err := parser.ParseFile(fset, filename, src, mode)
	if err != nil {
		return nil, err
	}
	return Decorate(fset, f).(*dst.File), nil
}

func ParseExpr(x string) (dst.Expr, error) {
	expr, err := parser.ParseExpr(x)
	if err != nil {
		return nil, err
	}
	return Decorate(nil, expr).(dst.Expr), nil
}

func ParseDir(fset *token.FileSet, path string, filter func(os.FileInfo) bool, mode parser.Mode) (map[string]*dst.Package, error) {
	pkgs, err := parser.ParseDir(fset, path, filter, mode)
	if err != nil {
		return nil, err
	}
	d := New()
	out := map[string]*dst.Package{}
	for k, v := range pkgs {
		out[k] = d.Decorate(fset, v).(*dst.Package)
	}
	return out, nil
}

func ParseExprFrom(fset *token.FileSet, filename string, src interface{}, mode parser.Mode) (dst.Expr, error) {
	expr, err := parser.ParseExprFrom(fset, filename, src, mode)
	if err != nil {
		return nil, err
	}
	return Decorate(fset, expr).(dst.Expr), nil
}

func Decorate(fset *token.FileSet, n ast.Node) dst.Node {
	return New().Decorate(fset, n)
}

func DecorateFile(fset *token.FileSet, f *ast.File) *dst.File {
	return New().Decorate(fset, f).(*dst.File)
}

func New() *Decorator {
	return &Decorator{
		Nodes:       map[ast.Node]dst.Node{},
		Scopes:      map[*ast.Scope]*dst.Scope{},
		Objects:     map[*ast.Object]*dst.Object{},
		decorations: map[ast.Node]map[string][]string{},
	}
}

type Decorator struct {
	Nodes       map[ast.Node]dst.Node
	Scopes      map[*ast.Scope]*dst.Scope
	Objects     map[*ast.Object]*dst.Object
	decorations map[ast.Node]map[string][]string
}

func (d *Decorator) Decorate(fset *token.FileSet, n ast.Node) dst.Node {
	fragger := &Fragger{}
	fragger.Fragment(fset, n)
	d.decorations = fragger.Link()
	return d.decorateNode(n)
}

type decorationInfo struct {
	name string
	decs []string
}

func (d *Decorator) decorateObject(o *ast.Object) *dst.Object {
	if o == nil {
		return nil
	}
	if do, ok := d.Objects[o]; ok {
		return do
	}
	/*
		// An Object describes a named language entity such as a package,
		// constant, type, variable, function (incl. methods), or label.
		//
		// The Data fields contains object-specific data:
		//
		//	Kind    Data type         Data value
		//	Pkg     *Scope            package scope
		//	Con     int               iota for the respective declaration
		//
		type Object struct {
			Kind ObjKind
			Name string      // declared name
			Decl interface{} // corresponding Field, XxxSpec, FuncDecl, LabeledStmt, AssignStmt, Scope; or nil
			Data interface{} // object-specific data; or nil
			Type interface{} // placeholder for type information; may be nil
		}
	*/

	out := &dst.Object{}
	d.Objects[o] = out
	out.Kind = dst.ObjKind(o.Kind)
	out.Name = o.Name

	switch decl := o.Decl.(type) {
	case *ast.Scope:
		out.Decl = d.decorateScope(decl)
	case ast.Node:
		out.Decl = d.decorateNode(decl)
	default:
		panic(fmt.Sprintf("o.Decl is %T", o.Decl))
	}

	// TODO: I believe Data is either a *Scope or an int. We will support both and panic if something else if found.
	switch data := o.Data.(type) {
	case int:
		out.Data = data
	case *ast.Scope:
		out.Data = d.decorateScope(data)
	case ast.Node:
		out.Data = d.decorateNode(data)
	case nil:
	default:
		panic(fmt.Sprintf("o.Data is %T", o.Data))
	}

	return out
}

func (d *Decorator) decorateScope(s *ast.Scope) *dst.Scope {
	if s == nil {
		return nil
	}
	if ds, ok := d.Scopes[s]; ok {
		return ds
	}
	/*
		// A Scope maintains the set of named language entities declared
		// in the scope and a link to the immediately surrounding (outer)
		// scope.
		//
		type Scope struct {
			Outer   *Scope
			Objects map[string]*Object
		}
	*/
	out := &dst.Scope{}

	d.Scopes[s] = out

	out.Outer = d.decorateScope(s.Outer)
	out.Objects = map[string]*dst.Object{}
	for k, v := range s.Objects {
		out.Objects[k] = d.decorateObject(v)
	}

	return out
}
