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
	defMakePath := fmt.Sprintf("var/rpkgm/main/%s/binary", pkg)
	cmd := exec.Command("make")
	cmd.Dir = defMakePath
	_, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to install '%s': %s", pkg, err)
		os.Exit(1)
	}
	log.Printf("Installation done for: %s.", pkg)
}

func uninstall(pkg string) {
	// Uninstall a package
	log.Printf("Uninstalling %s...", pkg)
	defMakePath := fmt.Sprintf("var/rpkgm/main/%s", pkg)
	cmd := exec.Command("make", "uninstall")
	cmd.Dir = defMakePath
	_, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to uninstall '%s': %s", pkg, err)
		os.Exit(1)
	}
	log.Printf("Uninstallation done for: %s.", pkg)
}
