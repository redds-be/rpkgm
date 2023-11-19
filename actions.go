package main

import (
	"fmt"
	"log"
	"os/exec"
)

func install(pkg string) {
	defMakePath := fmt.Sprintf("var/rpkgm/main/%s", pkg)
	cmd := exec.Command("make")
	cmd.Dir = defMakePath
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to install '%s' : %s", pkg, err)
	}
}

func uninstall(pkg string) {
	defMakePath := fmt.Sprintf("var/rpkgm/main/%s", pkg)
	cmd := exec.Command("make", "uninstall")
	cmd.Dir = defMakePath
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to uninstall '%s' : %s", pkg, err)
	}
}
