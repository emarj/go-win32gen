package gowin32gen

import (
	"strings"
	"text/template"
)

func (f Function) Generate() string {

	tmpl, err := template.New("func").Parse(
		`\\ Function {{.Name}} [DllImport "{{.DllImport}}"]
		func {{.Name}}({{- range $index,$p := .Params -}}{{- if $index -}}, {{- end -}}
			{{$p.Name}} {{$p.Type.Name}}
		   {{- end -}}) {{- .ReturnType.Name -}}
		`,
	)
	if err != nil {
		panic(err)
	}

	b := &strings.Builder{}

	err = tmpl.Execute(b, f)
	if err != nil {
		panic(err)
	}

	return b.String()
}
