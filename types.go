package gowin32gen

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type ApiFile struct {
	Constants      []Constant
	Types          []Type
	Functions      []json.RawMessage
	UnicodeAliases []json.RawMessage
}

const (
	BaseTypeKindNative    = "Native"
	BaseTypeKindApiRef    = "ApiRef"
	BaseTypeKindLPArray   = "LPArray"
	BaseTypeKindPointerTo = "PointerTo"
)

func ConvertBaseType(typeName string) (bool, string) {
	switch typeName {
	case "Int8", "Int16", "Int32", "Int64",
		"UInt8", "UInt16", "UInt32", "UInt64", "Float32", "Float64":
		return true, strings.ToLower(typeName)
	default:
		return false, ""
	}

}

type BaseType struct {
	Kind string
	Name *string `json:",omitempty"`
	// Kind = ApiRef
	BaseTypeApiRef
	// Kind = Array
	BaseTypeArray
	// Kind = LPArray
	BaseTypeLPArray
	// Kind = PointerTo, LPArray
	Child *BaseType `json:",omitempty"`
}

type BaseTypeArray struct {
	Shape *struct{ Size int } `json:",omitempty"`
}

type BaseTypeApiRef struct {
	TargetKind *string  `json:",omitempty"`
	Api        *string  `json:",omitempty"`
	Parents    []string `json:",omitempty"`
}

type BaseTypeLPArray struct {
	NullNullTerm    *bool `json:",omitempty"`
	CountConst      *int  `json:",omitempty"`
	CountParamIndex *int  `json:",omitempty"`
}

type Constant struct {
	Name      string
	Type      BaseType
	ValueType string
	//Value     ConstantValue
	Value json.RawMessage
	Attrs []string
}

func (c Constant) String() string {
	switch c.Type.Kind {
	case BaseTypeKindNative:
		isNative, nativeType := ConvertBaseType(*c.Type.Name)
		if isNative {
			return fmt.Sprintf("const %s %s = %s", c.Name, nativeType, c.Value)
		}
		log.Printf("WARNING: type %q flagged as native (Ref: %s)\n", *c.Type.Name, c.Name)
	case BaseTypeKindApiRef:
		return fmt.Sprintf("var %s %s = %s", c.Name, *c.Type.Name, c.Value)
	}
	return ""
}

// This is needed to handle both int values and PKs:
// "Value":15
// "Value":{"Fmtid":"1da5d803-d492-4edd-8c23-e0c0ffee7f0e","Pid":0}
type ConstantValue json.RawMessage

const (
	TypeKindEnum       = "Enum"
	TypeKindFP         = "FunctionPointer"
	TypeKindStruct     = "Struct"
	TypeKindCom        = "Com"
	TypeKindComClassID = "ComClassID"
)

type Type struct {
	Name          string
	Kind          string
	ValueType     string
	Architectures []string
	Platform      *string
	Guid          *string
	TypeEnum
	TypeStruct
	TypeFunctionPointer
	TypeCom
}

type TypeEnum struct {
	Flags       bool
	Values      []EnumValue
	IntegerBase string
}

type EnumValue struct {
	Name  string
	Value int
}

type TypeStruct struct {
	Scoped      bool
	Size        int
	PackingSize int
	Fields      []NamedBaseType
	NestedTypes []Type
}

type NamedBaseType struct {
	Name string
	Type BaseType
}

type TypeFunctionPointer struct {
	SetLastError bool
	Params       []FunParam
	ReturnType   BaseType
	ReturnAttrs  []string
}

type TypeCom struct {
	Interface BaseType
	Methods   []Type
}

type FunParam struct {
	Name  string
	Type  BaseType
	Attrs json.RawMessage
}

func (t Type) Generate() string {

	if t.Kind != TypeKindCom {
		return ""
	}

	w := StringWriter{}

	vtbl := t.Name + "VTable"
	iface := *t.Interface.Api + "." + *t.Interface.Name
	ifaceVtable := iface + "VTable"

	w.WriteLine(fmt.Sprintf("type %s struct{", t.Name))
	w.WriteLine(iface)
	w.WriteLine("}")

	w.WriteLine(fmt.Sprintf("type %s struct{\n", vtbl))
	w.WriteLine(ifaceVtable)
	for _, m := range t.Methods {
		w.WriteLine(m.Name + " uintptr")
	}
	w.WriteLine("}")

	w.WriteString("\n\n")

	w.WriteString(fmt.Sprintf(`
	func (v *%s) VTable() *%s {
		return (*%s)(unsafe.Pointer(v.RawVTable))
	}`, t.Name, vtbl, vtbl))

	w.WriteString("\n\n")

	for _, m := range t.Methods {
		w.WriteString(fmt.Sprintf(`
		func (v *%s) %s() %s{
			
		}`, t.Name, m.Name, *m.ReturnType.Name))
		w.WriteString("\n\n")
	}

	return w.String()

}

type StringWriter struct {
	strings.Builder
}

func (w *StringWriter) WriteLine(s string) error {
	_, err := w.WriteString(s + "\n")
	return err
}
