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

package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/redds-be/rpkgm/internal/database"
	"github.com/redds-be/rpkgm/internal/util"
	"github.com/spf13/cobra"
)

// Global vars for booleans acting as toggles.
var (
	showNotice  bool
	showLicense bool
	showInfo    bool
	showAll     bool
)

// copyright/warranty notice.
const notice = `rpgkm, redd's package manager.
Copyright (C) 2024 redd

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.`

// showCmd represents the show command.
var showCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "show",
	Short: "Show license related informations",
	Run: func(cmd *cobra.Command, args []string) {
		// show rpkgm's warranty notice
		if showNotice {
			util.Display(os.Stdout, false, "%s", notice)
			os.Exit(0)
		}

		// Connect to the database
		dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
		if err != nil {
			util.Display(os.Stderr, true, "rpkgm could connect to the repo's database. Error: %s", err)
			os.Exit(1)
		}

		if showAll {
			// Get every package info in the database
			allPkgInfo, err := dbAdapter.GetAllPkgInfo()
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could query the repo's database. Error: %s", err)
				os.Exit(1)
			}

			// For every package, display the info (changes depending on the installation status)
			for _, pkgInfo := range allPkgInfo {
				if pkgInfo.Installed {
					util.Display(
						os.Stdout, false,
						"%s [Installed (%s)]\t- %s\t- Repo's version: %s",
						pkgInfo.Name,
						pkgInfo.InstalledVersion,
						pkgInfo.Description,
						pkgInfo.RepoVersion,
					)
				} else {
					util.Display(os.Stdout, false, "%s [Not installed]\t- %s\t- Repo's version: %s", pkgInfo.Name, pkgInfo.Description, pkgInfo.RepoVersion)
				}
			}

			// Close the database
			err = dbAdapter.CloseDBConnection()
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not close the connection to the database. Error: %s", err)
				os.Exit(1)
			}

			os.Exit(0)
		}

		// The following operations require a package name, error if there isn't
		if name == "" {
			// Display help
			util.Display(os.Stderr, true, "Error: you need to specify a package with --name/-n first.")
			err := cmd.Help()
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not display the help message. Error: %s", err)
				os.Exit(1)
			}

			// Close the database
			err = dbAdapter.CloseDBConnection()
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not close the connection to the database. Error: %s", err)
				os.Exit(1)
			}
			os.Exit(1)
		}

		// Check if the package is in the repo, error this is not the case
		isInRepo, _ := dbAdapter.IsPkgInRepo(name)
		if !isInRepo {
			util.Display(os.Stderr, true, "The package: %s is not in repository.", name)

			// Close the database
			err = dbAdapter.CloseDBConnection()
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not close the connection to the database. Error: %s", err)
				os.Exit(1)
			}
			os.Exit(1)
		}

		if showLicense {
			// Find the build files for the given package, the license should be in there
			buildFilesDir, err = dbAdapter.GetPkgBuildFilesDir(name)
			if err != nil {
				util.Display(
					os.Stderr,
					true,
					"rpkgm could not find the build files (build files includes the license) for the given package. Error: %s",
					err,
				)
				os.Exit(1)
			}

			// Open the license file
			licenseFile, err := os.Open(fmt.Sprintf("%s/LICENSE", buildFilesDir))
			if err != nil {
				util.Display(
					os.Stderr, true,
					"rpkgm could not open or find the license file for the given package. Error: %s",
					err,
				)
				os.Exit(1)
			}

			// Defer the closing of the license file
			defer func(licenseFile *os.File) {
				err := licenseFile.Close()
				if err != nil {
					util.Display(
						os.Stderr, true,
						"rpkgm could not close the license file for the given package. Error: %s",
						err,
					)
					os.Exit(1)
				}
			}(licenseFile)

			// Read the license file
			licenseContent, err := io.ReadAll(licenseFile)
			if err != nil {
				util.Display(
					os.Stderr, true,
					"rpkgm could not read from the license file of the given package. Error: %s",
					err,
				)
				os.Exit(1)
			}

			// Print the license file's content
			util.Display(os.Stdout, false, "%v", string(licenseContent))
		}

		if showInfo {
			// Get the given package's general info
			pkgInfo, err := dbAdapter.GetPkgInfo(name)
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not get the given package's information. Error: %s", err)
				os.Exit(1)
			}

			// Display the general info (changes depending on the installation status)
			if pkgInfo.Installed {
				util.Display(
					os.Stdout, false,
					"%s [Installed (%s)]\t- %s\t- Repo's version: %s",
					pkgInfo.Name,
					pkgInfo.InstalledVersion,
					pkgInfo.Description,
					pkgInfo.RepoVersion,
				)
			} else {
				util.Display(os.Stdout, false, "%s [Not installed]\t- %s\t- Repo's version: %s", pkgInfo.Name, pkgInfo.Description, pkgInfo.RepoVersion)
			}
		}

		// Close the database connection
		err = dbAdapter.CloseDBConnection()
		if err != nil {
			util.Display(os.Stderr, true, "rpkgm could not close the connection to the database. Error: %s", err)
			os.Exit(1)
		}
	},
}

// init initializes the command-line arguments for cobra.
func init() { //nolint:gochecknoinits
	// Link to root (root = 'rpkgm', show = 'rpkgm show')
	rootCmd.AddCommand(showCmd)

	// --warranty or -w is used to print the GPLv3 warranty notice
	showCmd.Flags().BoolVarP(&showNotice, "warranty", "w", false, "Show warranty.")

	// Flag to select a package to show things
	showCmd.Flags().StringVarP(&name, "name", "n", "", "Package to select.")

	// Flag to show a given package's license
	showCmd.Flags().BoolVarP(&showLicense, "license", "l", false, "Show a given package's license.")

	// Flag to show a given package's general information.
	showCmd.Flags().BoolVarP(&showInfo, "info", "i", false, "Show a given package's general information.")

	// Flag to show every packages general information.
	showCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show every package's general information.")

	// Optional flag to specify repo database location
	showCmd.Flags().
		StringVarP(&repoDB, "repo", "r", "var/rpkgm/main.db", "Specify repo Database location (Defaults to /var/rpkgm/main.db).")
}
