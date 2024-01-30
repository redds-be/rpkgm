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
	"log"
	"os"
	"strings"

	"github.com/redds-be/rpkgm/internal/database"
	"github.com/redds-be/rpkgm/internal/util"
	"github.com/spf13/cobra"
)

// Global vars for information to add.
var (
	name          string
	description   string
	version       string
	buildFilesDir string
)

// addCmd represents the add command.
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add package to the main repo",
	Long:  `Add a package to the main repo (one at a time), using it's name and current version.'`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if user is root.
		util.CheckRoot()

		// Connect to the database
		dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
		if err != nil {
			util.Display(os.Stderr, "rpkgm could not connect to the database. Error: %s", err)
			os.Exit(1)
		}

		// Create the main repo table if it does not exist
		err = dbAdapter.CreatePkgTable()
		if err != nil {
			util.Display(os.Stderr, "rpkgm could not create the main repo table. Error: %s", err)
			os.Exit(1)
		}

		// Default value for buildFilesDir (doing it here instead of Flags() because I need 'name')
		if buildFilesDir == "" {
			buildFilesDir = fmt.Sprintf("var/rpkgm/main/%s", name)
		}

		// Remove any trailing /
		buildFilesDir = strings.TrimSuffix(buildFilesDir, "/")

		// Add the package to the main repo
		err = dbAdapter.AddToMainRepo(name, description, version, buildFilesDir)
		if err != nil {
			util.Display(os.Stderr, "rpkgm could not add the package to the repo. Error: %s", err)
			os.Exit(1)
		}

		// Close the database connection
		err = dbAdapter.CloseDBConnection()
		if err != nil {
			util.Display(os.Stderr, "rpkgm could not close the connection to the database. Error: %s", err)
			os.Exit(1)
		}
	},
}

// init initializes the command-line arguments for cobra.
func init() { //nolint:gochecknoinits
	// Link to root (root = 'rpkgm', add = 'rpkgm add')
	rootCmd.AddCommand(addCmd)

	// Package name
	addCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the package to add.")

	// Name is required
	err := addCmd.MarkFlagRequired("name")
	if err != nil {
		log.Fatal(err)
	}

	// Package version
	addCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the package to add.")

	// Version is required
	err = addCmd.MarkFlagRequired("version")
	if err != nil {
		log.Fatal(err)
	}

	// Optional flag to specify build files location
	addCmd.Flags().
		StringVarP(&buildFilesDir, "files", "f", "",
			"Build files location. (Defaults to /var/rpkgm/main/<pkgName>)")

	// Optional flag to give a description of a package
	addCmd.Flags().
		StringVarP(&description, "desc", "d", "[No description provided for this package.]", "Description of the package to add.")

	// Optional flag to specify repo database location
	addCmd.Flags().
		StringVarP(&repoDB, "repo", "r", "var/rpkgm/main.db", "Specify repo Database location (Defaults to /var/rpkgm/main.db).")
}
