package decorator

import (
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
		decorations: map[ast.Node]map[string][]string{},
	}
}

type Decorator struct {
	Nodes       map[ast.Node]dst.Node
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
