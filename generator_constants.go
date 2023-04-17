package gowin32gen

import (
	"fmt"
	"log"
)

func (c Constant) Generate() string {
	switch c.Type.Kind {
	case TypeKindNative:
		isNative, nativeType := ConvertType(*c.Type.Name)
		if isNative {
			return fmt.Sprintf("const %s %s = %s", c.Name, nativeType, c.Value)
		}
		log.Printf("WARNING: type %q flagged as native (Ref: %s)\n", *c.Type.Name, c.Name)
	case TypeKindApiRef:
		return fmt.Sprintf("var %s %s = %s", c.Name, *c.Type.Api+"."+*c.Type.Name, c.Value)
	}
	return ""
}
