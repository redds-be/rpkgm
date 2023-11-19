package main

import (
	"github.com/pborman/getopt/v2"
	"os"
)

func main() {
	// Temporary path
	path := "/var/log/rpkgm.log"

	err := getArgs()
	if err != nil {
		getopt.Usage()
		os.Exit(0)
	}

	logFile := handleLogs(path)
	defer closeLogs(logFile)

	installArg := getopt.GetValue("install")
	uninstallArg := getopt.GetValue("uninstall")

	if installArg != "" {
		install(installArg)
	} else if uninstallArg != "" {
		uninstall(uninstallArg)
	}
}
