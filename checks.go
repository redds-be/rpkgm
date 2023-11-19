package main

import (
	"fmt"
	"net/http"
	"os"
)

func checkRoot() {
	// Check if the user is root
	if os.Geteuid() != 0 {
		fmt.Println("Please run rpkgm as root!")
		os.Exit(1)
	}
}

func checkNetwork() {
	// Check the network connectivity
	_, err := http.Get("https://github.com")
	if err != nil {
		fmt.Println("You need to be connected to the internet to run rpkgm.")
		os.Exit(1)
	}
}

func doChecks() {
	// Driver for the checks
	checkRoot()
	checkNetwork()
}
