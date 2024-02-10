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
	"strings"

	"github.com/redds-be/rpkgm/internal/manage"
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
	archiveURL       string
	hash             string
	dependencies     []string
	remove           bool
)

// manageCmd represents the manage command.
var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Manage a package.",
	Long:  `Manage a given package using package's information in a given (or the default) repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if user is root.
		util.CheckRoot("Please run rpkgm manage as root.")

		// Literally everything here needs a package's name, error if there isn't one
		if name == "" {
			util.Display(
				os.Stderr,
				true,
				"You need to specify a package name with --name/-n before managing its information.",
			)
			err := cmd.Help()
			if err != nil {
				util.Display(
					os.Stderr,
					true,
					"rpkgm could not display the help message. Error: %s",
					err,
				)
				os.Exit(1)
			}
		}

		deps := ""
		if len(dependencies) > 0 {
			// Convert the dependencies list into a string
			deps = strings.Join(dependencies, " ")
		}

		// Decide what to do and do what is needed to do
		manage.Decide(
			repoDB,
			name,
			newName,
			newDesc,
			installedVersion,
			repoVersion,
			archiveURL,
			hash,
			deps,
			remove,
			markInstalled,
			markUninstalled,
		)
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
	manageCmd.Flags().
		StringVarP(&newDesc, "desc", "d", "", "Change the description of a given package.")

	// Flag to mark a package as installed
	manageCmd.Flags().
		BoolVarP(&markInstalled, "installed", "i", false, "Mark a given package as installed")

	// Flag to mark a package as not installed
	manageCmd.Flags().
		BoolVarP(&markUninstalled, "uninstalled", "u", false, "Mark a given package as not installed")

	// Mark installed and uninstalled as incompatible together
	manageCmd.MarkFlagsMutuallyExclusive("installed", "uninstalled")

	// Flag to change a package's installed version in the db
	manageCmd.Flags().
		StringVar(&installedVersion, "iv", "", "Change a given package's installed version in the database.")

	// Flag to change a package's repo version in the db
	manageCmd.Flags().
		StringVar(&repoVersion, "rv", "", "Change a given package's repo version in the database.")

	// Flag to change a package's archive URL
	manageCmd.Flags().
		StringVarP(&archiveURL, "archive", "a", "", "Change a given package's archive URL.")

	// Flag to change a package's archive's hash
	manageCmd.Flags().StringVar(&hash, "hash", "", "Change a given package's archive's hash.")

	// Flag to change a package's dependencies
	manageCmd.Flags().
		StringSliceVar(&dependencies, "deps", nil, "List of dependencires separated by commas for a given package.")

	// Flag to remove a package from the repo
	manageCmd.Flags().BoolVar(&remove, "rm", false, "Remove a given package from the repository.")

	// Optional flag to specify repo database location
	manageCmd.Flags().
		StringVarP(&repoDB, "repo", "r", "var/rpkgm/main/main.db", "Specify repo Database location.")
}
