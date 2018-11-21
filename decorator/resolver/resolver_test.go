package resolver_test

import (
	"testing"

	"path/filepath"

	"github.com/dave/dst/decorator/resolver"
)

func TestPackageResolver(t *testing.T) {
	type tc struct{ importPath, fromDir, expectName string }
	tests := []struct {
		skip, solo bool
		name       string
		resolve    func() (end func(), root string, r resolver.PackageResolver)
		cases      []tc
	}{
		{
			name: "resolver.Guess",
			resolve: func() (end func(), root string, r resolver.PackageResolver) {
				r = resolver.Guess{
					"a/b/c": "d",
				}
				return
			},
			cases: []tc{
				{"a/b/c", "", "d"},
				{"d/e/f", "", "f"},
			},
		},
		{
			name: "resolver.Map",
			resolve: func() (end func(), root string, r resolver.PackageResolver) {
				r = resolver.Map{
					"a/b/c": "d",
				}
				return
			},
			cases: []tc{
				{"a/b/c", "", "d"},
				{"d/e/f", "", ""},
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
				name, err := r.ResolvePackage(c.importPath, fromDir)
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
