package goast

import (
	"go/token"
	"testing"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

func TestGoAstDecoratorResolver(t *testing.T) {
	type tc struct{ id, expect string }
	tests := []struct {
		skip, solo bool
		name       string
		src        string
		cases      []tc
	}{
		{
			name: "simple",
			src: `package main

				import (
					"root/a"
				)

				func main(){
					a.A()
				}`,
			cases: []tc{
				{"A", "root/a"},
			},
		},
		{
			name: "shadow",
			src: `package main

				import (
					"root/a"
				)

				func main(a T){
					a.A()
				}

				type T struct{}
				func (T) A()`,
			cases: []tc{
				{"A", ""},
			},
		},
	}
	var solo bool
	for _, test := range tests {
		if test.solo {
			solo = true
			break
		}
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if solo && !test.solo {
				t.Skip()
			}
			if test.skip {
				t.Skip()
			}

			d := decorator.NewDecoratorWithImports(token.NewFileSet(), "main", New())

			f, err := d.Parse(test.src)
			if err != nil {
				panic(err)
			}

			nodes := map[string]string{}
			dst.Inspect(f, func(n dst.Node) bool {
				switch n := n.(type) {
				case *dst.Ident:
					nodes[n.Name] = n.Path
				}
				return true
			})

			for _, c := range test.cases {
				found, ok := nodes[c.id]
				if !ok {
					t.Errorf("node %s not found", c.id)
				}
				if found != c.expect {
					t.Errorf("expect %q, found %q", c.expect, found)
				}
			}

		})
	}
}
