package guess_test

import (
	"testing"

	"github.com/dave/dst/decorator/resolver"
	"github.com/dave/dst/decorator/resolver/guess"
)

func TestRestorerResolver(t *testing.T) {
	type tc struct{ importPath, expectName string }
	tests := []struct {
		skip, solo bool
		name       string
		resolve    func() (end func(), r resolver.RestorerResolver)
		cases      []tc
	}{
		{
			name: "guess.RestorerResolver",
			resolve: func() (end func(), r resolver.RestorerResolver) {
				r = guess.RestorerResolver{
					"a/b/c": "d",
				}
				return
			},
			cases: []tc{
				{"a/b/c", "d"},
				{"d/e/f", "f"},
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
				end, r := test.resolve()
				name, err := r.ResolvePackage(c.importPath)
				if end != nil {
					end() // delete temp dir if created
				}
				if err == resolver.ErrPackageNotFound {
					name = ""
				} else if err != nil {
					t.Errorf("error resolving path %s: %v", c.importPath, err)
				}
				if name != c.expectName {
					t.Errorf("package %s - expected %s, got %s", c.importPath, c.expectName, name)
				}
			}
		})
	}
}
