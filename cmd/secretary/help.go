package main

import (
	_ "embed"
	"flag"
	"fmt"
)

//go:embed help.txt
var helpText string

func printHelp() {
	fmt.Println(helpText)
	flag.PrintDefaults()
}
