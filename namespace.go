package gowin32gen

import (
	"fmt"
	"path"
	"strings"
)

var SourceFilesHaveExtension = true
var NamespaceSep = "."

func Filename2Pkg(source string) Pkg {
	base := path.Base(source)
	if SourceFilesHaveExtension {
		pos := strings.LastIndex(base, ".")
		if pos < 0 {
			panic(fmt.Sprintf("malformed input file %q", base))
		}
		base = base[:pos]
	}

	return PkgFromString(base)

}

type Pkg struct {
	Name  string
	Parts []string
}

func (p Pkg) String() string {
	if len(p.Parts) > 0 {
		return strings.Join([]string{strings.Join(p.Parts, NamespaceSep), p.Name}, NamespaceSep)
	} else {
		return p.Name
	}
}

func (p Pkg) Path() string {
	if len(p.Parts) > 0 {
		return strings.Join([]string{strings.Join(p.Parts, "/"), p.Name}, "/")
	} else {
		return p.Name
	}
}

func PkgFromString(s string) Pkg {
	return PkgFromSlice(strings.Split(s, "."))
}

func PkgFromSlice(s []string) Pkg {
	if len(s) == 0 {
		return Pkg{}
	}
	if len(s) == 1 {
		return Pkg{Name: s[0]}
	}
	return Pkg{Name: s[len(s)-1], Parts: s[:len(s)-1]}
}
