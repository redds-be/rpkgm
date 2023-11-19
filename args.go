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
	getopt.BoolLong("verbose", 'v', "Verbosity")
	getopt.ListLong("install", 'i', "Package name (for multiple packages, they need to be comma separated without spaces)")
	getopt.ListLong("uninstall", 'u', "Package name (for multiple packages, they need to be comma separated without spaces)")
	getopt.Parse()

	// Is -h (--help) called, if so, print the usage string (auto-generated from getopt)
	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	installArg := getopt.GetValue("install")
	uninstallArg := getopt.GetValue("uninstall")

	// Are -i (--install) and -u (--uninstall) both set ? A bit conflicting. Neither ? What is the point of calling rpkgm ?
	if installArg != "" && uninstallArg != "" {
		fmt.Println("The options -i (--install) and -u (--uninstall) are mutually exclusive.")
		return errors.New("-i and -u are mutually exclusive")
	} else if installArg == "" && uninstallArg == "" {
		fmt.Println("You should either use -i (--install) or -u (--uninstall).")
		return errors.New("neither -i nor -u are declared")
	}

	return nil
}
