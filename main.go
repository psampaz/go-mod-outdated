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
	flag.Parse()

	err := runner.Run(os.Stdin, os.Stdout, *withUpdate, *onlyDirect)

	if err != nil {
		log.Print(err)
	}
}
