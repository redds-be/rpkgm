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
	"os"

	"github.com/redds-be/rpkgm/internal/database"
	"github.com/redds-be/rpkgm/internal/util"
)

// pkgConf defines some of the general information about a package.
type pkgConf struct {
	pkgName string
	version string
	index   int
	total   int
	verbose bool
	keep    bool
	force   bool
}

// ask is used to ask the user to confirm before installing or uninstalling.
func ask(toInstall, toUninstall []string, dbAdapter database.Adapter, force bool) { //nolint:cyclop,funlen,gocognit
	var atLeastOne bool
	var choice string

	// if it's for installation, ask before installing, if it's for uninstallation, ask before uninstalling
	if len(toInstall) > 0 { //nolint:nestif
		for _, pkgName := range toInstall {
			// Check if the package is in the repo
			isInRepo, _ := dbAdapter.IsPkgInRepo(pkgName)
			if !isInRepo {
				util.Display(os.Stdout, "The package: %s is not in the repository. Skipping...", pkgName)

				continue
			}

			// Check if the package is already installed, if the user uses -f, we don't care
			isInstalled, _ := dbAdapter.IsInstalled(pkgName)
			if isInstalled && !force {
				util.Display(
					os.Stdout,
					"The package: %s is already installed (you can use --force/-f to force the installation). Skipping...",
					pkgName,
				)

				continue
			} else if isInstalled && force {
				util.Display(os.Stdout, "Install %s", pkgName)
				atLeastOne = true

				continue
			}

			// If every check is ok, we go on
			if isInRepo && !isInstalled {
				util.Display(os.Stdout, "Install %s", pkgName)
				atLeastOne = true
			}
		}

		// If there's at least one valid package on the install list, we ask before installing it.
		if atLeastOne {
			// Ask the user's confirmation
			fmt.Printf("Do you want to install these packages? [y/N] ") //nolint:forbidigo
			_, err := fmt.Scanln(&choice)
			if err != nil {
				// Close the database connection
				err = dbAdapter.CloseDBConnection()
				if err != nil {
					util.Display(os.Stderr, "rpkgm could not close the connection to the database. Error: %s", err)
					os.Exit(1)
				}
				os.Exit(0)
			}

			// Choice is yes, we go on
			if choice == "y" || choice == "Y" {
				return
			}
		} else {
			util.Display(os.Stderr, "No package selected for installation.")
		}

		// Close the database connection
		err := dbAdapter.CloseDBConnection()
		if err != nil {
			util.Display(os.Stderr, "rpkgm could not close the connection to the database. Error: %s", err)
			os.Exit(1)
		}
		os.Exit(0)
	} else if len(toUninstall) > 0 {
		for _, pkgName := range toUninstall {
			// Check if the package is in the repo
			isInRepo, _ := dbAdapter.IsPkgInRepo(pkgName)
			if !isInRepo {
				util.Display(os.Stdout, "The package: %s is not in the repository. Skipping...", pkgName)

				continue
			}

			// Check if the package is installed, if it isn't we skip it
			isInstalled, _ := dbAdapter.IsInstalled(pkgName)
			if !isInstalled {
				util.Display(
					os.Stdout,
					"The package: %s is not installed. Skipping...",
					pkgName,
				)

				continue
			}

			// If the package is in the repo and installed, we mark it for uninstallation
			if isInRepo && isInstalled {
				util.Display(os.Stdout, "Uninstall %s", pkgName)
				atLeastOne = true
			}
		}

		// If there's at least one valid package on the uninstall list, we ask before uninstalling it.
		if atLeastOne {
			// Ask the user's confirmation
			fmt.Printf("Do you want to uninstall these packages? [y/N] ") //nolint:forbidigo
			_, err := fmt.Scanln(&choice)
			if err != nil {
				// Close the database connection
				err = dbAdapter.CloseDBConnection()
				if err != nil {
					util.Display(os.Stderr, "rpkgm could not close the connection to the database. Error: %s", err)
					os.Exit(1)
				}
				os.Exit(0)
			}

			// If the choice is yes, we go on
			if choice == "y" || choice == "Y" {
				return
			}
		} else {
			util.Display(os.Stderr, "No package selected for uninstallation.")
		}

		// Close the database connection
		err := dbAdapter.CloseDBConnection()
		if err != nil {
			util.Display(os.Stderr, "rpkgm could not close the connection to the database. Error: %s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
}

func Manage( //nolint:funlen,cyclop,gocognit
	toInstall, toUninstall []string,
	verbose, keep, force, yes bool,
	repoDB string,
) {
	// Check if user is root
	util.CheckRoot()

	// Connect to the database
	dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
	if err != nil {
		util.Display(os.Stderr, "rpkgm could not connect to the database. Error: %s", err)
		os.Exit(1)
	}

	// Ask before installing/uninstalling unless --yes/-y is used
	if !yes {
		ask(toInstall, toUninstall, *dbAdapter, force)
	}

	// If there's a package marked for installation, we prepare for the installation,
	// if there's a package marked for uninstallation, we prepare for the uninstallation
	if len(toInstall) > 0 { //nolint:nolintlint,nestif
		for index, pkgName := range toInstall {
			// Since --yes/-y skips the first checks, we check if the package is in the repo
			if yes {
				isInRepo, _ := dbAdapter.IsPkgInRepo(pkgName)
				if !isInRepo {
					util.Display(os.Stderr, "The package named: %s is not in the repo.", pkgName)

					continue
				}
			}

			// Since --yes/-y skips the first checks, we check if the package is already installed
			if yes {
				// Check if the package is installed, skip if this the case
				isInstalled, err := dbAdapter.IsInstalled(pkgName)
				if err != nil && !force {
					util.Display(
						os.Stderr,
						"rpkgm could not determine if %s is installed, you can ignore this error with --force/-f. Error: %s",
						pkgName,
						err,
					)

					continue
				}

				// Skip unless the user uses --force/-f
				if isInstalled && !force {
					util.Display(
						os.Stderr,
						"The package named: %s is already marked as installed, to force install it, use --force/-f.",
						pkgName,
					)

					continue
				}
			}

			// Get the package's version
			version, err := dbAdapter.GetRepoVersion(pkgName)
			if err != nil && !force {
				util.Display(
					os.Stderr,
					"rpkgm could not determine %s's version in the repo. you can ignore this error with --force/-f. Error: %s",
					pkgName,
					err,
				)

				continue
			}

			// Put all the info into a struct to pass to the install func
			conf := pkgConf{
				pkgName: pkgName,
				version: version,
				index:   index + 1,
				total:   len(toInstall),
				verbose: verbose,
				keep:    keep,
				force:   force,
			}

			// Call the install func
			if err := conf.install(); err != nil {
				util.Display(os.Stderr, "%s", err)

				continue
			}

			// Mark the installed package as installed
			err = dbAdapter.MarkAsInstalled(pkgName)
			if err != nil {
				util.Display(
					os.Stderr,
					"rpkgm could not mark %s as installed although the package is, in fact installed. Error: %s",
					pkgName, err,
				)
			}

			// Set the package's installed version to the repo's version
			err = dbAdapter.SetInstalledVersion(pkgName, version)
			if err != nil {
				util.Display(
					os.Stderr,
					"rpkgm could not set %s's install version in the repo, although the package is, in fact installed. Error: %s",
					pkgName,
					err,
				)
			}
		}
	} else if len(toUninstall) > 0 {
		for index, pkgName := range toUninstall {
			// Since --yes/-y skips the first checks, we check if the package is in the repo and if it's installed
			if yes {
				// If the package is not in the repo, we skip it
				isInRepo, _ := dbAdapter.IsPkgInRepo(pkgName)
				if !isInRepo {
					util.Display(os.Stderr, "The package named: %s is not in the repo.", pkgName)

					continue
				}

				// If the package is not installed, we skip it
				isInstalled, err := dbAdapter.IsInstalled(pkgName)
				if err != nil {
					util.Display(os.Stderr, "rpkgm could not determine if %s is installed. Error: %s", pkgName, err)

					continue
				}

				if !isInstalled {
					util.Display(os.Stderr, "The package named: %s is not installed.", pkgName)

					continue
				}
			}

			// Put all the info into a struct to pass to the uninstall func
			conf := pkgConf{
				pkgName: pkgName,
				index:   index + 1,
				total:   len(toUninstall),
				verbose: verbose,
				keep:    keep,
			}

			// Call the uninstall func
			if err := conf.uninstall(); err != nil {
				util.Display(os.Stderr, "%s", err)

				continue
			}

			// Mark the package as not installed
			err = dbAdapter.MarkAsNotInstalled(pkgName)
			if err != nil {
				util.Display(os.Stderr, "rpkgm could not mark the package as uninstalled although the package is, in fact uninstalled. Error: %s", err)
			}

			// Remove the package's installed version from the db
			err = dbAdapter.SetInstalledVersion(pkgName, "")
			if err != nil {
				util.Display(os.Stderr, "rpkgm could not reset %s's install version in the repo, although the package is, in fact uninstalled. Error: %s", pkgName, err)
			}
		}
	}

	// Close the database connection
	err = dbAdapter.CloseDBConnection()
	if err != nil {
		util.Display(os.Stderr, "rpkgm could not close the connection to the database. Error: %s", err)
		os.Exit(1)
	}
}
