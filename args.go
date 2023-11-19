package main

import (
	"errors"
	"fmt"
	"github.com/pborman/getopt/v2"
)

func getArgs() error {
	optInstall := getopt.StringLong("install", 'i', "", "Package name")
	optUninstall := getopt.StringLong("uninstall", 'u', "", "Package name")
	optHelp := getopt.BoolLong("help", 'h', "Help")
	getopt.Parse()

	if *optHelp {
		return errors.New("an unexpected error occurred")
	}

	if *optInstall != "" && *optUninstall != "" {
		fmt.Println("The options -i (--install) and -u (--uninstall) are mutually exclusive.")
		return errors.New("-i and -u are mutually exclusive")
	}

	if *optInstall != "" {
		return nil
	} else if *optUninstall != "" {
		return nil
	} else {
		return errors.New("an unexpected error occurred")
	}
}
