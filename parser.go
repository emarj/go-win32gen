package gowin32gen

import (
	"encoding/json"
	"log"
	"os"
)

type APIFile struct {
	Pkg        Pkg
	DLLImports []string
	Data       File
}

type APIRef struct {
	Name string
	Pkg  Pkg
}

type Parser struct {
	parsedFiles    map[string]APIFile
	requiredAPIRef []APIRef
}

func NewParser() *Parser {
	return &Parser{parsedFiles: map[string]APIFile{}}
}

func (p *Parser) ParseFile(filename string) error {
	source, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer source.Close()

	pkg := Filename2Pkg(filename)

	log.Printf("Parsing file %q (%q)", filename, pkg)

	data := File{}
	decoder := json.NewDecoder(source)
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}

	p.parsedFiles[pkg.String()] = APIFile{Pkg: pkg, Data: data}

	return err
}
