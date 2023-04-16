package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/emarj/gowin32gen"
)

const defaultOutputFolder = "./.api"

var outputFolder string

func main() {

	flag.StringVar(&outputFolder, "o", defaultOutputFolder, "output folder")

	flag.Parse()

	if _, err := os.Stat(outputFolder); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if flag.NArg() == 0 {
		fmt.Println("usage: gowin32gen <source> -o <dest>")
		flag.PrintDefaults()
		fmt.Println(os.Args)
		os.Exit(1)
	}

	filelist := flag.Args()

	for _, filename := range filelist {
		err := ProcessFile(filename)
		if err != nil {
			fmt.Printf("error %q: %s\n", filename, err)
		}
	}
}

func ProcessFile(filename string) error {
	source, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer source.Close()

	destname := FilenameTransform(filename)
	err = os.MkdirAll(path.Dir(destname), 0777)
	if err != nil {
		return err
	}

	dest, err := os.Create(destname)
	if err != nil {
		return err
	}
	defer dest.Close()

	return gowin32gen.Generate(source, dest)

}

func FilenameTransform(source string) string {
	pieces := strings.Split(path.Base(source), ".")
	if len(pieces) == 1 {
		return source
	}
	dest := outputFolder

	for i := 0; i < len(pieces)-1; i++ {
		dest = path.Join(dest, pieces[i])
	}
	dest += ".go"

	return dest

}
