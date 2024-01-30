//    rpkgm, redd's package manager.
//    Copyright (C) 2024 redd
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with this program.  If not, see <https://www.gnu.org/licenses/>.

package pkg

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/redds-be/rpkgm/internal/util"
)

// install executes the installation script's commands.
func (pc pkgConf) install() error { //nolint:funlen,cyclop
	// Set the location of the rbuild
	rbuild := fmt.Sprintf("var/rpkgm/main/%s/rbuild", pc.pkgName)
	// Set the working directory, if -k, set it to /usr/src/rpkgm
	var workdir string
	if pc.keep {
		workdir = fmt.Sprintf("/tmp/usr/src/rpkgm/%s", pc.pkgName)
	} else {
		workdir = fmt.Sprintf("/tmp/rpkgm/%s", pc.pkgName)
	}

	// Resolve dependencies
	printDeps := fmt.Sprintf("source %s ; printDeps", rbuild)
	depsOut, err := exec.Command("/usr/bin/env", "bash", "-c", printDeps).CombinedOutput()
	if err != nil {
		log.Printf("rpkgm could not resolve dependencies for: %s\n", pc.pkgName)

		return err
	}

	// Inform of the dependencies
	if string(depsOut) != "" {
		util.Display(
			os.Stdout,
			"The package %s, marks [%s] as dependencies, rpkgm not support dependency resolution and"+
				" installation yet, you will need to install them manually.",
			pc.pkgName,
			string(depsOut),
		)
	}

	// Inform of the downloading
	util.Display(
		os.Stdout,
		"Downloading (%s%d%s of %s%d%s) %s%s=%s%s",
		util.By,
		pc.index,
		util.Rc,
		util.By,
		pc.total,
		util.Rc,
		util.Bg,
		pc.pkgName,
		pc.version,
		util.Rc,
	)

	// Download the package's archive
	download := fmt.Sprintf("source %s ; download %s", rbuild, workdir)
	dlOut, err := exec.Command("/usr/bin/env", "bash", "-c", download).CombinedOutput()
	if err != nil {
		util.Display(os.Stderr, "rpkgm could not download the archive for: %s", pc.pkgName)

		return err
	}

	// Display the output
	if pc.verbose && string(dlOut) != "" {
		log.Println(string(dlOut))
	}

	// Do we verify the archive?
	if !pc.force {
		// Verify the package's archive
		verify := fmt.Sprintf("source %s ; verify %s", rbuild, workdir)
		verOut, err := exec.Command("/usr/bin/env", "bash", "-c", verify).CombinedOutput()
		if err != nil {
			util.Display(os.Stderr, "rpkgm could not verify the archive for: %s\n%s", pc.pkgName, string(verOut))
			util.Display(os.Stderr, "You can disable the archive verification by re-running using --force/-f.")

			return err
		}

		// Display the output
		if pc.verbose && string(verOut) != "" {
			util.Display(os.Stdout, string(verOut))
		}
	}

	// Inform of the extraction
	util.Display(
		os.Stdout,
		"Extracting (%s%d%s of %s%d%s) %s%s=%s%s",
		util.By,
		pc.index,
		util.Rc,
		util.By,
		pc.total,
		util.Rc,
		util.Bg,
		pc.pkgName,
		pc.version,
		util.Rc,
	)

	// Extract the package's archive
	extract := fmt.Sprintf("source %s ; extract %s", rbuild, workdir)
	exOut, err := exec.Command("/usr/bin/env", "bash", "-c", extract).CombinedOutput()
	if err != nil {
		util.Display(os.Stderr, "rpkgm could not extract the archive for: %s", pc.pkgName)

		return err
	}

	// Display the output
	if pc.verbose && string(exOut) != "" {
		log.Println(string(exOut))
	}

	makefile := fmt.Sprintf("var/rpkgm/main/%s/Makefile", pc.pkgName)

	// Copy rpkgm's makefile into source dir.
	_, err = exec.Command("/usr/bin/cp", makefile, workdir).CombinedOutput()
	if err != nil {
		util.Display(os.Stderr, "rpkgm could not copy the makefile for: %s", pc.pkgName)

		return err
	}

	// Inform of the installing
	util.Display(
		os.Stdout,
		"Installing (%s%d%s of %s%d%s) %s%s=%s%s",
		util.By,
		pc.index,
		util.Rc,
		util.By,
		pc.total,
		util.Rc,
		util.Bg,
		pc.pkgName,
		pc.version,
		util.Rc,
	)

	// Install the package
	install := fmt.Sprintf("source %s ; install %s", rbuild, workdir)
	inOut, err := exec.Command("/usr/bin/env", "bash", "-c", install).CombinedOutput()
	if err != nil {
		util.Display(os.Stderr, "rpkgm could not install: %s", pc.pkgName)

		return err
	}

	// Display the output
	if pc.verbose && string(inOut) != "" {
		util.Display(os.Stdout, string(inOut))
	}

	// If we don't keep the source, remove it
	if !pc.keep {
		// Inform of the cleaning
		util.Display(
			os.Stdout,
			"Cleaning (%s%d%s of %s%d%s) %s%s=%s%s",
			util.By,
			pc.index,
			util.Rc,
			util.By,
			pc.total,
			util.Rc,
			util.Bg,
			pc.pkgName,
			pc.version,
			util.Rc,
		)

		// Clean working directory
		err = os.RemoveAll(workdir)
		if err != nil {
			util.Display(os.Stderr, "rpkgm could not clean the working directory for: %s", pc.pkgName)

			return err
		}
	}

	return nil
}
