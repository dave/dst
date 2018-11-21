package gobuild_test

import (
	"testing"

	"path/filepath"

	"github.com/dave/dst/decorator/resolver"
	"github.com/dave/dst/decorator/resolver/gobuild"
	"github.com/dave/dst/dstutil/dummy"
)

func TestPackageResolver(t *testing.T) {
	type tc struct{ importPath, fromDir, expectName string }
	tests := []struct {
		skip, solo bool
		name       string
		resolve    func() (end func(), root string, r *gobuild.PackageResolver)
		cases      []tc
	}{
		{
			name: "gobuild.Resolver",
			resolve: func() (end func(), root string, r *gobuild.PackageResolver) {
				src := dummy.Dir{
					"main1": dummy.Dir{
						"vendor":   dummy.Dir{"a": dummy.Dir{"a.go": dummy.Src("package a1 \n\n func A(){}")}},
						"main1.go": dummy.Src("package main \n\n func main(){}"),
					},
					"main2": dummy.Dir{
						"main2.go": dummy.Src("package main \n\n func main(){}"),
					},
					"a": dummy.Dir{"a.go": dummy.Src("package a2 \n\n func A(){}")},
				}
				r = &gobuild.PackageResolver{Context: dummy.BuildContext(src)}
				root = "/gopath/src"
				return
			},
			cases: []tc{
				{"a", "/main1", "a1"},
				{"a", "/main2", "a2"},
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
