package gowin32gen_test

import (
	"testing"

	"github.com/emarj/gowin32gen"
)

func TestFilename2Pkg(t *testing.T) {

	tests := map[string]gowin32gen.Pkg{
		"/ad/aasd/ada/asd/this.is.a.test.json": {Name: "test", Parts: []string{"this", "is", "a"}},
	}

	for input, expected := range tests {
		got := gowin32gen.Filename2Pkg(input)
		if got.String() != expected.String() {
			t.Errorf("%s != %s", got, expected)
		}
	}

}

func TestPkgPath(t *testing.T) {

	tests := map[string]gowin32gen.Pkg{
		"this/is/a/test": {Name: "test", Parts: []string{"this", "is", "a"}},
	}

	for expected, pkg := range tests {
		got := pkg.Path()
		if got != expected {
			t.Errorf("%s != %s", got, expected)
		}
	}

}
