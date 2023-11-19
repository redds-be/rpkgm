package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"os"
)

func main() {
	pkg, optType, err := getArgs()
	if err != nil {
		getopt.Usage()
		os.Exit(0)
	}

	if optType == "install" {
		fmt.Printf("Installing %s...\n", *pkg)
	} else if optType != "uninstall" {
		fmt.Printf("Uninstalling %s...\n", *pkg)
	}
}
