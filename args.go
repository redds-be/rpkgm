package main

import (
	"errors"
	"fmt"
	"github.com/pborman/getopt/v2"
)

func getArgs() (*string, string, error) {
	optInstall := getopt.StringLong("install", 'i', "", "Package name")
	optUninstall := getopt.StringLong("uninstall", 'u', "", "Package name")
	optHelp := getopt.BoolLong("help", 'h', "Help")
	getopt.Parse()

	if *optHelp {
		return nil, "", errors.New("an unexpected error occurred")
	}

	if *optInstall != "" && *optUninstall != "" {
		fmt.Println("The options -i (--install) and -u (--uninstall) are mutually exclusive.")
		return nil, "", errors.New("-i and -u are mutually exclusive")
	}

	if *optInstall != "" {
		return optInstall, "install", nil
	} else if *optUninstall != "" {
		return optUninstall, "uninstall", nil
	} else {
		return nil, "", errors.New("an unexpected error occurred")
	}
}
