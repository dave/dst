package gopackages_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dave/dst/decorator/resolver"
	"github.com/dave/dst/decorator/resolver/gopackages"
	"github.com/dave/dst/dstutil/dummy"
)

func TestPackageResolver(t *testing.T) {
	type tc struct{ importPath, fromDir, expectName string }
	tests := []struct {
		skip, solo bool
		name       string
		resolve    func() (end func(), root string, r *gopackages.PackageResolver)
		cases      []tc
	}{
		{
			name: "gopackages.Resolver",
			resolve: func() (end func(), root string, r *gopackages.PackageResolver) {
				src := dummy.Dir{
					"main":   dummy.Dir{"main.go": dummy.Src("package main \n\n func main(){}")},
					"foo":    dummy.Dir{"foo.go": dummy.Src("package foo \n\n func A(){}")},
					"go.mod": dummy.Src("module root"),
				}
				root = dummy.TempDir(src)
				end = func() { os.RemoveAll(root) }
				r = &gopackages.PackageResolver{}
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
				if err == resolver.PackageNotFoundError {
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
