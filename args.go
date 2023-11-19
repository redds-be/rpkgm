package main

import (
	"errors"
	"fmt"
	"github.com/pborman/getopt/v2"
)

func getArgs() error {
	optHelp := getopt.BoolLong("help", 'h', "Help")
	_ = getopt.BoolLong("verbose", 'v', "Verbosity")
	optInstall := getopt.StringLong("install", 'i', "", "Package name")
	optUninstall := getopt.StringLong("uninstall", 'u', "", "Package name")
	getopt.Parse()

	if *optHelp {
		return errors.New("an unexpected error occurred")
	}

	if *optInstall != "" && *optUninstall != "" {
		fmt.Println("The options -i (--install) and -u (--uninstall) are mutually exclusive.")
		return errors.New("-i and -u are mutually exclusive")
	} else if *optInstall == "" && *optUninstall == "" {
		fmt.Println("You should either use -i (--install) or -u (--uninstall).")
		return errors.New("neither -i nor -u are declared")
	}

	return nil
}
