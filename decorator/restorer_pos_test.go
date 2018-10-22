package decorator

import (
	"testing"
)

func TestPosTests(t *testing.T) {
	testPackageRestoresCorrectly(t, "github.com/dave/dst/gendst/data")
}
