package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	bf "github.com/russross/blackfriday/v2"
	"github.com/spf13/afero"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mdex <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), wd)

	markdownFile, err := afero.ReadFile(fs, filename)
	if err != nil {
		log.Fatal(err)
	}

	markdown := bf.New(
		bf.WithExtensions(bf.Extensions(bf.CodeBlock)),
	)

	walkFunc := func(node *bf.Node, entering bool) bf.WalkStatus {
		if entering && node.Type == bf.CodeBlock {
			info := string(node.CodeBlockData.Info)

			if filename, ok := strings.CutPrefix(info, "go "); ok {
				if len(filename) == 0 {
					return bf.GoToNext
				}

				fs.MkdirAll(filepath.Dir(filename), 0o755)
				afero.WriteFile(fs, filename, node.Literal, 0o644)

				fmt.Println(filename)
			}
		}

		return bf.GoToNext
	}

	node := markdown.Parse(markdownFile)
	node.Walk(walkFunc)
}
