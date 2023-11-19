package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"os"
)

func main() {
	// Parse the config
	conf, err := getConf()
	if err != nil {
		fmt.Printf("Coudln't parse the configuration file %s", err)
		os.Exit(1)
	}

	// Parse the args
	err = getArgs()
	if err != nil {
		getopt.Usage()
		os.Exit(0)
	}

	// Write the logs to a file and to stdout if verbosity is set
	logFile := handleLogs(conf.logFile, conf.verbose)
	defer closeLogs(logFile)

	// Set the values for install/uninstall args
	installArg := getopt.GetValue("install")
	uninstallArg := getopt.GetValue("uninstall")

	// Are we installing or uninstalling ?
	if installArg != "" {
		// We are installing
		install(installArg)
	} else if uninstallArg != "" {
		// We are uninstalling
		uninstall(uninstallArg)
	}
}
