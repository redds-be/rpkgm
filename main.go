package main

import (
	"github.com/pborman/getopt/v2"
	"log"
	"os"
)

func main() {
	// Temporary path
	path := "/var/log/rpkgm.log"

	logFile := handleLogs(path)
	defer closeLogs(logFile)
	log.Println("Starting rpkgm...")

	err := getArgs()
	if err != nil {
		getopt.Usage()
		os.Exit(0)
	}

	installArg := getopt.GetValue("install")
	uninstallArg := getopt.GetValue("uninstall")

	if installArg != "" {
		install(installArg)
	} else if uninstallArg != "" {
		uninstall(uninstallArg)
	}
}
