package gowin32gen

import (
	"io/fs"
	"os"
	"path"
	"strings"
)

var PermissionMask int = 0777

type Generator struct {
	*Parser
	OutFolder string
}

func NewGenerator(outFolder string) *Generator {
	return &Generator{OutFolder: outFolder, Parser: NewParser()}
}

func (g *Generator) Generate() error {
	for k, file := range g.Parser.parsedFiles {
		output := GenerateFile(file)
		err := WriteFile(g.OutFolder, PkgFromString(k), output)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteFile(folder string, pkg Pkg, data []byte) error {
	fullpath := path.Join(folder, path.Join(pkg.Parts...), pkg.Name)
	err := os.MkdirAll(fullpath, fs.FileMode(PermissionMask))
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(fullpath, pkg.Name+".go"), data, fs.FileMode(PermissionMask))
	if err != nil {
		return err
	}

	return nil
}

func GenerateFile(af APIFile) []byte {

	b := strings.Builder{}

	b.WriteString("package " + af.Pkg.Name + "\n\n")

	//_, err = dest.Write([]byte(fmt.Sprintf("%#v", data)))
	for _, c := range af.Data.Constants {
		b.WriteString(c.Generate())
		b.WriteByte('\n')
	}

	for _, t := range af.Data.Types {
		b.WriteString(t.Generate())
		b.WriteByte('\n')
	}

	for _, f := range af.Data.Functions {
		b.WriteString(f.Generate())
		b.WriteByte('\n')
	}

	return []byte(b.String())
}
