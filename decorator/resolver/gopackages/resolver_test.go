package gopackages_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dave/dst/decorator/resolver"
	"github.com/dave/dst/decorator/resolver/gopackages"
)

func TestRestorerResolver(t *testing.T) {
	type tc struct{ importPath, fromDir, expectName string }
	tests := []struct {
		skip, solo bool
		name       string
		resolve    func() (end func(), root string, r *gopackages.RestorerResolver)
		cases      []tc
	}{
		{
			name: "gopackages.Resolver",
			resolve: func() (end func(), root string, r *gopackages.RestorerResolver) {
				src := map[string]string{
					"main/main.go": "package main \n\n func main(){}",
					"foo/foo.go":   "package foo \n\n func A(){}",
					"go.mod":       "module root",
				}
				var err error
				root, err = tempDir(src)
				if err != nil {
					t.Fatal(err)
				}
				end = func() { os.RemoveAll(root) }
				r = &gopackages.RestorerResolver{}
				return
			},
			cases: []tc{
				{"root/foo", "/main", "foo"},
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
			for _, c := range test.cases {
				end, root, r := test.resolve()
				fromDir := filepath.Join(root, c.fromDir)
				r.Dir = fromDir
				name, err := r.ResolvePackage(c.importPath)
				if end != nil {
					end() // delete temp dir if created
				}
				if err == resolver.ErrPackageNotFound {
					name = ""
				} else if err != nil {
					t.Errorf("error resolving path %s from dir %s: %v", c.importPath, fromDir, err)
				}
				if name != c.expectName {
					t.Errorf("package %s, dir %s - expected %s, got %s", c.importPath, c.fromDir, c.expectName, name)
				}
			}
		})
	}
}
