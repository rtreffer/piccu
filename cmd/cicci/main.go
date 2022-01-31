package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rtreffer/piccu/pkg/cicci"
)

func main() {
	showHelp := flag.Bool("help", false, "displays a help text")
	flag.BoolVar(showHelp, "h", false, "displays a help text")

	flag.Parse()

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	// 1. find files

	files, err := cicci.CollectFiles(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't find files to merge:", err)
		os.Exit(1)
	}
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "no files found, refusing to create empty archive")
		os.Exit(2)
	}

	// 2. load / expand templates

	expanded, err := files.LoadAndExpand(nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't load/expand files:", err)
		os.Exit(3)
	}

	// 3. validate

	errors := expanded.Validate()
	for _, err := range errors {
		fmt.Fprintln(os.Stderr, "WARNING:", err)
	}

	// 4. generate multipart archive

	archive, err := cicci.CreateMultipartArchive(expanded)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't create multipart archive:", err)
		os.Exit(4)
	}

	// 5. write multipart archive to stdout

	fmt.Println(archive)
}
