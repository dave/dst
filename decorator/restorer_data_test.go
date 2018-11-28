package decorator

import "testing"

func TestData(t *testing.T) {
	testPackageRestoresCorrectlyWithImports(t, "github.com/dave/dst/gendst/data")
}
