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
	"encoding/json"
	"fmt"
	"io"
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
	importFile    string
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

		// If the user imports a file, add the packages described in it
		if importFile != "" {
			// Open the file to import
			jsonPkgFile, err := os.Open(importFile)
			if err != nil {
				util.Display(os.Stderr, "rpkgm couldn't open the json file. Error: %s", err)
				os.Exit(1)
			}

			// Read the file's content
			contentInbytes, err := io.ReadAll(jsonPkgFile)
			if err != nil {
				util.Display(
					os.Stderr,
					"rpkgm couldn't read the json file's content. Error: %s",
					err,
				)
				os.Exit(1)
			}

			// Initialize the Packages struct
			var pkgs database.Packages

			// Read the json content of the file
			err = json.Unmarshal(contentInbytes, &pkgs)
			if err != nil {
				util.Display(os.Stderr, "rpkgm couldn't read the json file. Error: %s", err)
				os.Exit(1)
			}

			// for every package in the json file, add it to the repo
			for index := 0; index < len(pkgs.Packages); index++ {
				// If there isn't a description, give one by default
				if pkgs.Packages[index].Description == "" {
					pkgs.Packages[index].Description = "[No description provided for this package.]"
				}

				// If there isn't a build files dir, give one by default
				if pkgs.Packages[index].BuildFilesDir == "" {
					pkgs.Packages[index].BuildFilesDir = fmt.Sprintf(
						"var/rpkgm/main/%s",
						pkgs.Packages[index].Name,
					)
				}

				// Add the package to the repo
				err = dbAdapter.AddToMainRepo(
					pkgs.Packages[index].Name,
					pkgs.Packages[index].Description,
					pkgs.Packages[index].Version,
					pkgs.Packages[index].BuildFilesDir,
				)
				if err != nil {
					util.Display(
						os.Stderr,
						"rpkgm was unable to add %s to the repo. Error: %s",
						pkgs.Packages[index].Name,
						err,
					)
				}
			}

			// Close the json file
			err = jsonPkgFile.Close()
			if err != nil {
				util.Display(os.Stderr, "rpgkm couln't close the json file. Error: %s", err)
				os.Exit(1)
			}
		}

		// If there's at least a name and a version, add the package
		if name != "" && version != "" {
			// Default value for buildFilesDir (doing it here instead of Flags() because I need 'name')
			if buildFilesDir == "" {
				buildFilesDir = fmt.Sprintf("var/rpkgm/main/%s", name)
			}

			// Remove any trailing /
			buildFilesDir = strings.TrimSuffix(buildFilesDir, "/")

			// Add the package to the main repo
			err = dbAdapter.AddToMainRepo(name, description, version, buildFilesDir)
			if err != nil {
				util.Display(
					os.Stderr,
					"rpkgm could not add the package to the repo. Error: %s",
					err,
				)
				os.Exit(1)
			}
		}

		// Close the database connection
		err = dbAdapter.CloseDBConnection()
		if err != nil {
			util.Display(
				os.Stderr,
				"rpkgm could not close the connection to the database. Error: %s",
				err,
			)
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

	// Package version
	addCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the package to add.")

	// If name is used, then version should be used too
	addCmd.MarkFlagsRequiredTogether("name", "version")

	// Optional flag to specify build files location
	addCmd.Flags().
		StringVarP(&buildFilesDir, "files", "f", "",
			"Build files location. (Defaults to /var/rpkgm/main/<pkgName>)")

	// Optional flag to give a description of a package
	addCmd.Flags().
		StringVarP(&description, "desc", "d", "[No description provided for this package.]", "Description of the package to add.")

	// Flag to import a json file containing the record to add to the repo's db
	addCmd.Flags().StringVarP(&importFile, "import", "i", "", "JSON file to import to the repo.")

	// Optional flag to specify repo database location
	addCmd.Flags().
		StringVarP(&repoDB, "repo", "r", "var/rpkgm/main.db", "Specify repo Database location (Defaults to /var/rpkgm/main.db).")
}
