package main

import (
	_ "embed"
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/sprig"
)

//go:embed help.txt
var helpText string

func printHelp() {
	fmt.Println(helpText)
	funcs := sprig.TxtFuncMap()
	keys := make([]string, 0, len(funcs))
	for k := range funcs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println("supported tempalte functions:")
	start := keys[0][0]
	p := 0
	for i, k := range keys {
		if k[0] == start {
			continue
		}
		fmt.Println("  |", strings.Join(keys[p:i], ","))
		p = i
		start = keys[i][0]
	}
	fmt.Println("  |", strings.Join(keys[p:], ","))
	fmt.Println("see https://masterminds.github.io/sprig/ for full function documentation")
	fmt.Println()
	flag.PrintDefaults()
}
