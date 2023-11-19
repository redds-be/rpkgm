package main

import (
	"errors"
	"fmt"
	"github.com/pborman/getopt/v2"
	"log"
)

func getArgs() error {
	optInstall := getopt.StringLong("install", 'i', "", "Package name")
	optUninstall := getopt.StringLong("uninstall", 'u', "", "Package name")
	optHelp := getopt.BoolLong("help", 'h', "Help")
	log.Println("Parsing options...")
	getopt.Parse()

	if *optHelp {
		log.Println("Printing help.")
		return errors.New("an unexpected error occurred")
	}

	if *optInstall != "" && *optUninstall != "" {
		log.Println("Options -i (--install) and -u (--uninstall) are mutually exclusive.")
		fmt.Println("The options -i (--install) and -u (--uninstall) are mutually exclusive.")
		return errors.New("-i and -u are mutually exclusive")
	}

	return nil
}
