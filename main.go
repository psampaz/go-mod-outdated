// Package main is the entry point of the go-mod-outdated tool
package main

import (
	"flag"
	"log"
	"os"

	"github.com/psampaz/go-mod-outdated/internal/runner"
)

func main() {

	withUpdate := flag.Bool("update", false, "List only modules with updates")
	onlyDirect := flag.Bool("direct", false, "List only direct modules")
	exitNonZero := flag.Bool("ci", false, "Non-zero exit code when at least one outdated dependency was found")
	flag.Parse()

	err := runner.Run(os.Stdin, os.Stdout, *withUpdate, *onlyDirect, *exitNonZero)

	if err != nil {
		log.Print(err)
	}
}
