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

// uninstall executes the Makefile's uninstallation commands.
func (pc pkgConf) uninstall() error { //nolint:funlen
	// Set the location of the rbuild
	rbuild := fmt.Sprintf("var/rpkgm/main/%s/rbuild", pc.pkgName)

	// Get package's version
	versionCmd := fmt.Sprintf("source %s ; echo -n $version", rbuild)
	version, err := exec.Command("/usr/bin/env", "bash", "-c", versionCmd).CombinedOutput()
	if err != nil {
		log.Printf("rpkgm could not find the version for: %s\n", pc.pkgName)

		return err
	}

	// Set the location (directory) of the makefile
	makefileDir := fmt.Sprintf("var/rpkgm/main/%s", pc.pkgName)

	// Inform of the uninstalling
	util.Display(
		os.Stdout,
		"Uninstalling (%s%d%s of %s%d%s) %s%s=%s%s",
		util.By,
		pc.index,
		util.Rc,
		util.By,
		pc.total,
		util.Rc,
		util.Bg,
		pc.pkgName,
		version,
		util.Rc,
	)

	// Uninstall the package
	uninstall := fmt.Sprintf("source %s ; uninstall %s", rbuild, makefileDir)
	unOut, err := exec.Command("/usr/bin/env", "bash", "-c", uninstall).CombinedOutput()
	if err != nil {
		util.Display(os.Stderr, "rpkgm could not uninstall: %s", pc.pkgName)

		return err
	}

	// Display the output
	if pc.verbose && string(unOut) != "" {
		log.Println(string(unOut))
	}

	// If we don't keep the source, remove it
	if !pc.keep {
		workdir := fmt.Sprintf("/tmp/usr/src/rpkgm/%s", pc.pkgName)
		if _, err := os.Stat(workdir); !os.IsNotExist(err) {
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
				version,
				util.Rc,
			)

			// Clean working directory
			err = os.RemoveAll(workdir)
			if err != nil {
				util.Display(os.Stderr, "rpkgm could not clean the working directory for: %s", pc.pkgName)

				return err
			}
		}
	}

	return nil
}
