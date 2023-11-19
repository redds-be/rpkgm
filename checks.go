package main

import (
	"fmt"
	"net/http"
	"os"
)

func earlyChecks() {
	// Driver for the checks
	if os.Geteuid() != 0 {
		fmt.Println("Please run rpkgm as root!")
		os.Exit(1)
	}
	_, err := http.Get("https://github.com")
	if err != nil {
		fmt.Println("You need to be connected to the internet to run rpkgm.")
		os.Exit(1)
	}
}
