package main

import (
	"fmt"
	"github.com/pborman/getopt/v2"
	"log"
	"os"
	"strings"
)

func main() {
	// Checking if rpkgm is running as root and testing network connectivity
	earlyChecks()

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
		// if -a (--ask) is set, we ask for confirmation before continuing
		if getopt.GetValue("ask") == "true" {
			log.Printf("Asking the user for confirmation before installing: %s", installArg)
			var confirmation string
			fmt.Printf("The following packages are marked for installation : %s\n", installArg)
			fmt.Printf("Do you want to proceed with those changes: (y/N) ")
			_, err := fmt.Scanln(&confirmation)
			if err != nil {
				log.Println("Quitting...")
				os.Exit(0)
			}
			if confirmation == "y" || confirmation == "Y" {
			} else {
				log.Println("Quitting...")
				os.Exit(0)
			}
		}
		pkgList := strings.Split(installArg, ",")
		total := len(pkgList)
		for count, pkg := range pkgList {
			install(pkg, count+1, total)
		}
	}
	if uninstallArg != "" {
		// We are uninstalling
		// if -a (--ask) is set, we ask for confirmation before continuing
		if getopt.GetValue("ask") == "true" {
			log.Printf("Asking the user for confirmation before uninstalling: %s", uninstallArg)
			var confirmation string
			fmt.Printf("The following packages are marked for uninstallation : %s\n", uninstallArg)
			fmt.Printf("Do you want to proceed with those changes: (y/N) ")
			_, err := fmt.Scanln(&confirmation)
			if err != nil {
				log.Println("Quitting...")
				os.Exit(0)
			}
			if confirmation == "y" || confirmation == "Y" {
			} else {
				log.Println("Quitting...")
				os.Exit(0)
			}
		}
		pkgList := strings.Split(uninstallArg, ",")
		total := len(pkgList)
		for count, pkg := range pkgList {
			uninstall(pkg, count+1, total)
		}
	}
}
