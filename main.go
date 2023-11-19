package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"os"
)

func main() {
	err := getArgs()
	if err != nil {
		getopt.Usage()
		os.Exit(0)
	}

	installArg := getopt.GetValue("install")
	uninstallArg := getopt.GetValue("uninstall")

	if installArg != "" {
		fmt.Printf("Installing %s...\n", installArg)
		install(installArg)
	} else if uninstallArg != "" {
		fmt.Printf("Uninstalling %s...\n", uninstallArg)
		uninstall(uninstallArg)
	}
}
