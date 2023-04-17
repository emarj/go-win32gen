package main

import (
	"flag"
	"fmt"
	"os"

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

	generator := gowin32gen.NewGenerator(outputFolder)

	for _, filename := range filelist {
		err := generator.ParseFile(filename)
		if err != nil {
			fmt.Printf("error %q: %s\n", filename, err)
		}
	}

	err := generator.Generate()
	if err != nil {
		fmt.Println(err)
	}
}
