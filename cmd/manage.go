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
	"os"

	"github.com/redds-be/rpkgm/internal/database"
	"github.com/redds-be/rpkgm/internal/util"
	"github.com/spf13/cobra"
)

var (
	newName          string
	newDesc          string
	markInstalled    bool
	markUninstalled  bool
	installedVersion string
	repoVersion      string
	remove           bool
)

// manageCmd represents the manage command.
var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Manage a package.",
	Long:  `Manage a given package using package's information in a given (or the default) repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if user is root.
		util.CheckRoot()

		// Literally everything here needs a package's name, error if there isn't
		if name == "" {
			util.Display(os.Stderr, true, "Error: you need to specify a package with --name/-n first.")
			err := cmd.Help()
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not display the help message. Error: %s", err)
				os.Exit(1)
			}
		}

		// Connect to the database
		dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
		if err != nil {
			util.Display(os.Stderr, true, "rpkgm could not connect to the database. Error: %s", err)
			os.Exit(1)
		}

		// Defer the closing of the database connection
		defer func(dbAdapter *database.Adapter) {
			err := dbAdapter.CloseDBConnection()
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not close the connection to the database. Error: %s", err)
				os.Exit(1)
			}
		}(dbAdapter)

		// Check if the given package is in the repo (forcing the close if the db connection since I use os.Exit)
		isInRepo, _ := dbAdapter.IsPkgInRepo(name)
		if !isInRepo {
			util.Display(os.Stderr, true, "The package: %s is not in the repository.", name)

			// Close the database connection
			err := dbAdapter.CloseDBConnection()
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not close the connection to the database. Error: %s", err)
				os.Exit(1)
			}
			os.Exit(1)
		}

		// Remove the given package
		if remove {
			err = dbAdapter.RemovePackage(name)
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not delete the package from the repository. Error: %s", err)
				os.Exit(1)
			}

			return
		}

		// Rename the given package
		if newName != "" {
			err = dbAdapter.RenamePackage(name, newName)
			if err != nil {
				util.Display(
					os.Stderr,
					true,
					"rpkgm could not rename the package in the repo's database. Error: %s",
					err,
				)
				os.Exit(1)
			}
		}

		// Change the given package's description
		if newDesc != "" {
			err = dbAdapter.ChangePkgDesc(name, newDesc)
			if err != nil {
				util.Display(
					os.Stderr, true,
					"rpkgm could not change the package's description in the repo's database. Error: %s",
					err,
				)
				os.Exit(1)
			}
		}

		// Mark the given package as installed
		if markInstalled {
			err = dbAdapter.MarkAsInstalled(name)
			if err != nil {
				util.Display(
					os.Stderr, true,
					"rpkgm could not mark the package as installed in the repo's database. Error: %s",
					err,
				)
				os.Exit(1)
			}
		}

		// Mark the given package as not installed
		if markUninstalled {
			err = dbAdapter.MarkAsNotInstalled(name)
			if err != nil {
				util.Display(
					os.Stderr, true,
					"rpkgm could not mark the package as not installed in the repo's database. Error: %s",
					err,
				)
				os.Exit(1)
			}
		}

		// Change or set the given package's installed version
		if installedVersion != "" {
			err = dbAdapter.SetInstalledVersion(name, installedVersion)
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not set the package's installed version. Error: %s", err)
				os.Exit(1)
			}
		}

		// Change or set the given package's repo version
		if repoVersion != "" {
			err = dbAdapter.UpdateRepoVersion(name, repoVersion)
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not set the package's repo version. Error: %s", err)
				os.Exit(1)
			}
		}
	},
}

// init initializes the command-line arguments for cobra.
func init() { //nolint:gochecknoinits
	// Link to root (root = 'rpkgm', manage = 'rpkgm manage')
	rootCmd.AddCommand(manageCmd)

	// Flag for the name of a package
	manageCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the package to manage.")

	// Flag for the new name of a package to be renamed
	manageCmd.Flags().StringVar(&newName, "ren", "", "Rename a given package.")

	// Flag for the new description to give to a package
	manageCmd.Flags().StringVarP(&newDesc, "desc", "d", "", "Change the description of a given package.")

	// Flag to mark a package as installed
	manageCmd.Flags().BoolVarP(&markInstalled, "installed", "i", false, "Mark a given package as installed")

	// Flag to mark a package as not installed
	manageCmd.Flags().BoolVarP(&markUninstalled, "uninstalled", "u", false, "Mark a given package as not installed")

	// Mark installed and uninstalled as incompatible together
	manageCmd.MarkFlagsMutuallyExclusive("installed", "uninstalled")

	// Flag to change a package's installed version in the db
	manageCmd.Flags().
		StringVar(&installedVersion, "iv", "", "Change a given package's installed version in the database.")

	// Flag to change a package's repo version in the db
	manageCmd.Flags().StringVar(&repoVersion, "rv", "", "Change a given package's repo version in the database.")

	// Flag to remove a package from the repo
	manageCmd.Flags().BoolVar(&remove, "rm", false, "Remove a given package from the repository.")

	// Optional flag to specify repo database location
	manageCmd.Flags().
		StringVarP(&repoDB, "repo", "r", "var/rpkgm/main.db", "Specify repo Database location (Defaults to /var/rpkgm/main.db).")
}
