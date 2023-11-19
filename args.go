package main

import (
	"errors"
	"fmt"
	"github.com/pborman/getopt/v2"
	"os"
)

func getArgs() error {
	// Parse config
	optHelp := getopt.BoolLong("help", 'h', "Help")
	_ = getopt.BoolLong("verbose", 'v', "Verbosity")
	optInstall := getopt.StringLong("install", 'i', "", "Package name")
	optUninstall := getopt.StringLong("uninstall", 'u', "", "Package name")
	getopt.Parse()

	// Is -h (--help) called, if so, print the usage string (auto-generated from getopt)
	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	// Are -i (--install) and -u (--uninstall) both set ? A bit conflicting. Neither ? What is the point of calling rpkgm ?
	if *optInstall != "" && *optUninstall != "" {
		fmt.Println("The options -i (--install) and -u (--uninstall) are mutually exclusive.")
		return errors.New("-i and -u are mutually exclusive")
	} else if *optInstall == "" && *optUninstall == "" {
		fmt.Println("You should either use -i (--install) or -u (--uninstall).")
		return errors.New("neither -i nor -u are declared")
	}

	return nil
}
