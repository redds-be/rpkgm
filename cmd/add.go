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
	"github.com/redds-be/rpkgm/internal/add"
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

		// Decide what to do and do what is needed to do
		add.Decide(repoDB, name, description, version, buildFilesDir, importFile)
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
		StringVarP(&buildFilesDir, "files", "f", "", "Build files location.")

	// Optional flag to give a description of a package
	addCmd.Flags().
		StringVarP(&description, "desc", "d", "[No description provided for this package.]", "Description of the package to add.")

	// Flag to import a json file containing the record to add to the repo's db
	addCmd.Flags().StringVarP(&importFile, "import", "i", "", "JSON file to import to the repo.")

	// Optional flag to specify repo database location
	addCmd.Flags().
		StringVarP(&repoDB, "repo", "r", "var/rpkgm/main/main.db", "Specify repo Database location.")
}
