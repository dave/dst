package decorator

import (
	"testing"
)

func TestLoad(t *testing.T) {
	testPackageRestoresCorrectlyWithImports(t, "github.com/dave/dst/gendst/data")
}
