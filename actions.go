package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func install(pkg string) {
	// Install a package
	log.Printf("Installing %s...", pkg)
	// Temporary "hard" coded path for packages
	defMakePath := fmt.Sprintf("var/rpkgm/main/%s/binary", pkg)

	// Inefficient way to give an info to a user
	depinfo := exec.Command("make", "depinfo")
	depinfo.Dir = defMakePath
	stdout, err := depinfo.Output()
	if err != nil {
		log.Printf("Failed to get dependency informations for '%s': %s", pkg, err)
		fmt.Printf("Failed to get dependency informations for '%s': %s", pkg, err)
	}
	fmt.Println(string(stdout))

	// Build and install the package
	cmd := exec.Command("make")
	cmd.Dir = defMakePath
	stdout, err = cmd.Output()
	if err != nil {
		log.Printf("Failed to install '%s': %s", pkg, err)
		fmt.Printf("Failed to install '%s': %s", pkg, err)
		os.Exit(1)
	}
	log.Println(string(stdout))
	log.Printf("Installation done for: %s.", pkg)
}

func uninstall(pkg string) {
	// Uninstall a package
	log.Printf("Uninstalling %s...", pkg)
	defMakePath := fmt.Sprintf("var/rpkgm/main/%s/binary", pkg)

	// Uninstall the package
	cmd := exec.Command("make", "uninstall")
	cmd.Dir = defMakePath
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to uninstall '%s': %s", pkg, err)
		fmt.Printf("Failed to uninstall '%s': %s", pkg, err)
		os.Exit(1)
	}
	log.Println(string(stdout))
	log.Printf("Uninstallation done for: %s.", pkg)
}
