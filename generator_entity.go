package gowin32gen

import (
	"fmt"
	"log"
	"strings"
	"text/template"
)

func (e Entity) Generate() string {

	switch e.Kind {
	case EntityKindEnum:
		return GenerateEntityEnum(e)
	case EntityKindNativeTypedef:
		return GenerateEntityNativeTypedef(e)
	/*case EntityKindStruct:

	case EntityKindFP:
	case EntityKindCom:
	case EntityKindComClassID: */
	default:
		log.Printf("INFO: Generator for entity of kind %q not implemented\n", e.Kind)
	}

	return ""

}

func GenerateEntityEnum(e Entity) string {
	ok, baseType := ConvertType(e.IntegerBase)
	if !ok {
		panic(fmt.Sprintf("enum %q has non native type %q", e.Name, e.IntegerBase))
	}
	//TODO: find another way to pass this into the template
	e.IntegerBase = baseType

	tmpl, err := template.New("enum").Parse(
		`\\ Enum {{.Name}}
		const (
			{{- $base := .IntegerBase -}}
			{{ range $val := .Values }}
     		{{$val.Name}} {{$base}} = {{$val.Value}}
			{{- end -}}
		)
		`,
	)
	if err != nil {
		panic(err)
	}

	b := &strings.Builder{}

	err = tmpl.Execute(b, e)
	if err != nil {
		panic(err)
	}

	return b.String()
}

func GenerateEntityNativeTypedef(e Entity) string {
	baseType := ""
	ok := false
	if e.Def.Kind == TypeKindNative {
		ok, baseType = ConvertType(*e.Def.Name)
		if !ok {
			panic(fmt.Sprintf("nativetypedef %q has non native type %q", e.Name, *e.Def.Name))
		}
	} else if e.Def.Kind == TypeKindPointerTo {
		ok, baseType = ConvertType(*e.Def.Child.Name)
		if !ok {
			panic(fmt.Sprintf("nativetypedef %q does not point to a native type %q", e.Name, *e.Def.Child.Name))
		}
	} else {
		log.Printf("WARNING: NativeTypedef with kind %q", e.Def.Kind)
		return ""
	}

	//TODO: find another way to pass this into the template
	e.Def.Name = &baseType

	tmpl, err := template.New("enum").Parse(
		`\\ NativeTypedef {{.Name}}
		type {{ .Name }} {{.Def.Name}}
		`,
	)
	if err != nil {
		panic(err)
	}

	b := &strings.Builder{}

	err = tmpl.Execute(b, e)
	if err != nil {
		panic(err)
	}

	return b.String()
}

func GenerateEntityCOM(e Entity) string {
	w := StringWriter{}

	vtbl := e.Name + "VTable"
	/* iface := *t.Interface.Api + "." + *t.Interface.Name
	ifaceVtable := iface + "VTable"
	*/
	w.WriteLine(fmt.Sprintf("type %s struct{", e.Name))
	//w.WriteLine(iface)
	w.WriteLine("}")

	w.WriteLine(fmt.Sprintf("type %s struct{\n", vtbl))
	//w.WriteLine(ifaceVtable)
	for _, m := range e.Methods {
		w.WriteLine(m.Name + " uintptr")
	}
	w.WriteLine("}")

	/* 	w.WriteString("\n\n")

	   	w.WriteString(fmt.Sprintf(`
	   	func (v *%s) VTable() *%s {
	   		return (*%s)(unsafe.Pointer(v.RawVTable))
	   	}`, t.Name, vtbl, vtbl)) */

	w.WriteString("\n\n")

	for _, m := range e.Methods {
		_, tn := ConvertType(*m.ReturnType.Name)
		w.WriteString(fmt.Sprintf(`
		func (v *%s) %s() %s{
			
		}`, e.Name, m.Name, tn))
		w.WriteString("\n\n")
	}

	return w.String()
}
