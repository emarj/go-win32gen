package gowin32gen

import (
	"encoding/json"
	"strings"
)

type File struct {
	Constants      []Constant
	Types          []Entity
	Functions      []Function
	UnicodeAliases []json.RawMessage
}

const (
	TypeKindNative    = "Native"
	TypeKindApiRef    = "ApiRef"
	TypeKindLPArray   = "LPArray"
	TypeKindPointerTo = "PointerTo"
)

func ConvertType(typeName string) (bool, string) {
	isNative := false
	outName := typeName

	switch typeName {
	case "Int8", "Int16", "Int32", "Int64",
		"UInt8", "UInt16", "UInt32", "UInt64", "Byte",
		"Float32", "Float64",
		"UIntPtr":
		isNative = true
		outName = strings.ToLower(typeName)
	case "Char":
		isNative = true
		outName = "byte"
	case "IntPtr":
		//TODO: we implicitly assume x64 arch. Infact this should be int32 in x86
		isNative = true
		outName = "int64"
	}

	return isNative, outName
}

type Type struct {
	Kind string
	Name *string `json:",omitempty"`
	// Kind = ApiRef
	TypeApiRef
	// Kind = Array
	TypeArray
	// Kind = LPArray
	TypeLPArray
	// Kind = PointerTo, LPArray
	Child *Type `json:",omitempty"`
}

type TypeArray struct {
	Shape *struct{ Size int } `json:",omitempty"`
}

type TypeApiRef struct {
	TargetKind *string  `json:",omitempty"`
	Api        *string  `json:",omitempty"`
	Parents    []string `json:",omitempty"`
}

type TypeLPArray struct {
	NullNullTerm    *bool `json:",omitempty"`
	CountConst      *int  `json:",omitempty"`
	CountParamIndex *int  `json:",omitempty"`
}

type Constant struct {
	Name      string
	Type      Type
	ValueType string
	Value     json.RawMessage
	Attrs     []string
}

// This is needed to handle both int values and PKs:
// "Value":15
// "Value":{"Fmtid":"1da5d803-d492-4edd-8c23-e0c0ffee7f0e","Pid":0}
type ConstantValue json.RawMessage

const (
	EntityKindEnum          = "Enum"
	EntityKindFP            = "FunctionPointer"
	EntityKindStruct        = "Struct"
	EntityKindCom           = "Com"
	EntityKindComClassID    = "ComClassID"
	EntityKindNativeTypedef = "NativeTypedef"
)

type Entity struct {
	Name          string
	Kind          string
	ValueType     string
	Architectures []string
	Platform      *string
	Guid          *string
	EntityEnum
	EntityStruct
	EntityFunctionPointer
	EntityCom
	EntityNativeTypedef
}

type EntityEnum struct {
	Flags  bool
	Values []struct {
		Name  string
		Value int
	}
	IntegerBase string
}

type EntityStruct struct {
	Scoped      bool
	Size        int
	PackingSize int
	Fields      []struct {
		Name  string
		Type  Type
		Attrs []string
	}
	NestedTypes []Entity
}

type EntityNativeTypedef struct {
	AlsoUsableFor *string
	Def           *Type
	FreeFunc      *string
}

type EntityFunctionPointer struct {
	SetLastError bool
	Params       []FunParam
	ReturnType   Type
	ReturnAttrs  []string
}

type EntityCom struct {
	Interface Type
	Methods   []Entity
}

type FunParam struct {
	Name  string
	Type  Type
	Attrs json.RawMessage
}

type StringWriter struct {
	strings.Builder
}

func (w *StringWriter) WriteLine(s string) error {
	_, err := w.WriteString(s + "\n")
	return err
}

type Function struct {
	Architectures []string
	Platform      *string
	Name          string
	DllImport     string
	EntityFunctionPointer
}
