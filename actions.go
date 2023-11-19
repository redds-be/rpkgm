package main

import (
	"fmt"
	"github.com/gookit/color"
	"log"
	"os"
	"os/exec"
)

func install(pkg string, count int, total int) {
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

	// Download the package
	color.Printf(">>> Downloading (<yellow>%d</> of <yellow>%d</>) <green>%s</>\n", count, total, pkg)
	cmd := exec.Command("make", "dl")
	cmd.Dir = defMakePath
	stdout, err = cmd.Output()
	if err != nil {
		log.Printf("Failed to download '%s': %s", pkg, err)
		fmt.Printf("Failed to download '%s': %s", pkg, err)
		os.Exit(1)
	}
	log.Println(string(stdout))

	// Extract the package
	color.Printf(">>> Extracting (<yellow>%d</> of <yellow>%d</>) <green>%s</>\n", count, total, pkg)
	cmd = exec.Command("make", "extract")
	cmd.Dir = defMakePath
	stdout, err = cmd.Output()
	if err != nil {
		log.Printf("Failed to extract '%s': %s", pkg, err)
		fmt.Printf("Failed to extract '%s': %s", pkg, err)
		os.Exit(1)
	}
	log.Println(string(stdout))

	// Install the package
	color.Printf(">>> Installing (<yellow>%d</> of <yellow>%d</>) <green>%s</>\n", count, total, pkg)
	cmd = exec.Command("make", "install")
	cmd.Dir = defMakePath
	stdout, err = cmd.Output()
	if err != nil {
		log.Printf("Failed to install '%s': %s", pkg, err)
		fmt.Printf("Failed to install '%s': %s", pkg, err)
		os.Exit(1)
	}
	log.Println(string(stdout))

	// Clean the build directory
	color.Printf(">>> Cleaning (<yellow>%d</> of <yellow>%d</>) <green>%s</>\n", count, total, pkg)
	cmd = exec.Command("make", "clean")
	cmd.Dir = defMakePath
	stdout, err = cmd.Output()
	if err != nil {
		log.Printf("Failed to clean '%s': %s", pkg, err)
		fmt.Printf("Failed to clean '%s': %s", pkg, err)
		os.Exit(1)
	}
	log.Println(string(stdout))
	color.Printf(">>> Finished installing (<yellow>%d</> of <yellow>%d</>) <green>%s</>\n", count, total, pkg)
	log.Printf("Installation done for: %s.", pkg)
}

func uninstall(pkg string, count int, total int) {
	// Uninstall a package
	log.Printf("Uninstalling %s...", pkg)
	defMakePath := fmt.Sprintf("var/rpkgm/main/%s/binary", pkg)

	// Uninstall the package
	color.Printf(">>> Uninstalling (<yellow>%d</> of <yellow>%d</>) <green>%s</>\n", count, total, pkg)
	cmd := exec.Command("make", "uninstall")
	cmd.Dir = defMakePath
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to uninstall '%s': %s", pkg, err)
		fmt.Printf("Failed to uninstall '%s': %s", pkg, err)
		os.Exit(1)
	}
	log.Println(string(stdout))
	color.Printf(">>> Finished uninstalling (<yellow>%d</> of <yellow>%d</>) <green>%s</>\n", count, total, pkg)
	log.Printf("Uninstallation done for: %s.", pkg)
}
