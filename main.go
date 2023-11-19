package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"os"
)

func main() {
	// Temporary path
	conf, err := getConf()
	if err != nil {
		fmt.Printf("Coudln't parse the configuration file %s", err)
		os.Exit(1)
	}

	err = getArgs()
	if err != nil {
		getopt.Usage()
		os.Exit(0)
	}

	logFile := handleLogs(conf.logFile, conf.verbose)
	defer closeLogs(logFile)

	installArg := getopt.GetValue("install")
	uninstallArg := getopt.GetValue("uninstall")

	if installArg != "" {
		install(installArg)
	} else if uninstallArg != "" {
		uninstall(uninstallArg)
	}
}
