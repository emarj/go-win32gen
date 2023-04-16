package gowin32gen_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/emarj/gowin32gen"
)

var baseTypesExamples []string = []string{
	`{"Kind":"Array","Shape":{"Size":32},"Child":{"Kind":"ApiRef","Name":"CHAR","TargetKind":"Default","Api":"Foundation"}}`,
}

func TestBaseTypes(t *testing.T) {
	var bt gowin32gen.BaseType
	for _, ex := range baseTypesExamples {

		err := json.Unmarshal([]byte(ex), &bt)
		if err != nil {
			t.Fatal(err, ex)
			continue
		}

		output, err := json.Marshal(bt)
		if err != nil {
			t.Fatal(err, ex)
			continue
		}

		if !bytes.Equal(output, []byte(ex)) {
			t.Fatalf("%s != %s", ex, output)
		}

	}
}
