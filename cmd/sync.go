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
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/redds-be/rpkgm/internal/database"
	"github.com/redds-be/rpkgm/internal/util"
	"github.com/spf13/cobra"
)

var remote string

// syncCmd represents the sync command.
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the main or a specified repo.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if user is root.
		util.CheckRoot()

		if importFile == "" {
			destDir := fmt.Sprintf("var/rpkgm/%s", name)

			err := os.MkdirAll(destDir, os.ModePerm)
			if err != nil {
				util.Display(
					os.Stderr,
					true,
					"rpkgm could not create the destination directory for the build files.",
				)
				os.Exit(1)
			}

			archive := fmt.Sprintf("var/rpkgm/%s/%s.tar.gz", name, name)

			err = util.Download(
				archive,
				fmt.Sprintf("https://%s/raw/main/%s.tar.gz", remote, name),
			)
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not download the build files.")
				os.Exit(1)
			}

			importFile = fmt.Sprintf("var/rpkgm/%s/repo.json", name)

			err = util.Download(
				importFile,
				fmt.Sprintf("https://%s/raw/main/repo.json", remote),
			)
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not download the JSON file of the repo.")
				os.Exit(1)
			}

			err = util.Untar(destDir, archive)
			if err != nil {
				util.Display(
					os.Stderr,
					true,
					"rpkgm could not untar the repo's archive. Error: %s",
					err,
				)
				os.Exit(1)
			}
		}

		// Open the file to import
		jsonPkgFile, err := os.Open(importFile)
		if err != nil {
			util.Display(os.Stderr, true, "rpkgm couldn't open the json file. Error: %s", err)
			os.Exit(1)
		}

		// Read the file's content
		contentInbytes, err := io.ReadAll(jsonPkgFile)
		if err != nil {
			util.Display(
				os.Stderr, true,
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
			util.Display(os.Stderr, true, "rpkgm couldn't read the json file. Error: %s", err)
			os.Exit(1)
		}

		if _, err := os.Stat(repoDB); errors.Is(err, os.ErrNotExist) {
			// Connect to the database
			dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
			if err != nil {
				util.Display(
					os.Stderr,
					true,
					"rpkgm could not connect to the database. Error: %s",
					err,
				)
				os.Exit(1)
			}

			// Create the packages table if it does not exist
			err = dbAdapter.CreatePkgTable()
			if err != nil {
				util.Display(
					os.Stderr,
					true,
					"rpkgm could not create the packages table in the repo. Error: %s",
					err,
				)
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
				err = dbAdapter.AddToRepo(
					pkgs.Packages[index].Name,
					pkgs.Packages[index].Description,
					pkgs.Packages[index].Version,
					pkgs.Packages[index].BuildFilesDir,
				)
				if err != nil {
					util.Display(
						os.Stderr, true,
						"rpkgm was unable to add %s to the repo. Error: %s",
						pkgs.Packages[index].Name,
						err,
					)
				}
			}

			// Close the database connection
			err = dbAdapter.CloseDBConnection()
			if err != nil {
				util.Display(
					os.Stderr, true,
					"rpkgm could not close the connection to the database. Error: %s",
					err,
				)
				os.Exit(1)
			}
		} else {
			// Connect to the database
			dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
			if err != nil {
				util.Display(os.Stderr, true, "rpkgm could not connect to the database. Error: %s", err)
				os.Exit(1)
			}
			// for every package in the json file, update their record
			for index := 0; index < len(pkgs.Packages); index++ {
				// Since UPDATE will set null strings as null in the database,
				// we need to get the old information to avoid deleting some information if the JSON field is null.
				pkgInfo, err := dbAdapter.GetPkgInfo(pkgs.Packages[index].Name)
				if err != nil {
					util.Display(
						os.Stderr,
						true,
						"rpkgm was unable to get the information for the package %s to compare against new information. Error: %s",
						pkgs.Packages[index].Name,
						err,
					)

					continue
				}

				// If there isn't a description, give one by default
				if pkgs.Packages[index].Description == "" {
					pkgs.Packages[index].Description = pkgInfo.Description
				}

				// If there isn't a build files dir, give one by default
				if pkgs.Packages[index].BuildFilesDir == "" {
					pkgs.Packages[index].BuildFilesDir = pkgInfo.BuildFilesDir
				}

				if pkgs.Packages[index].Version == "" {
					pkgs.Packages[index].Version = pkgInfo.RepoVersion
				}

				// Add the package to the repo
				err = dbAdapter.SyncRepo(
					pkgs.Packages[index].Name,
					pkgs.Packages[index].Description,
					pkgs.Packages[index].Version,
					pkgs.Packages[index].BuildFilesDir,
				)
				if err != nil {
					util.Display(
						os.Stderr, true,
						"rpkgm was unable to update %s in the repo. Error: %s",
						pkgs.Packages[index].Name,
						err,
					)
				}
			}

			// Close the database connection
			err = dbAdapter.CloseDBConnection()
			if err != nil {
				util.Display(
					os.Stderr, true,
					"rpkgm could not close the connection to the database. Error: %s",
					err,
				)
				os.Exit(1)
			}
		}

		// Close the json file
		err = jsonPkgFile.Close()
		if err != nil {
			util.Display(os.Stderr, true, "rpgkm couln't close the json file. Error: %s", err)
			os.Exit(1)
		}
	},
}

// init initializes the command-line arguments for cobra.
func init() { //nolint:gochecknoinits
	// Link to root (root = 'rpkgm', sync = 'rpkgm sync')
	rootCmd.AddCommand(syncCmd)

	// Flag for a JSON file to sync with the repo's db
	syncCmd.Flags().
		StringVarP(&importFile, "import", "i", "", "JSON file to sync with the repo.")

	// Optional flag to specify repo database location
	syncCmd.Flags().
		StringVar(&repoDB, "repo", "var/rpkgm/main/main.db", "Specify repo Database location.")

	// Optional flag to specify a remote GitHub repo (only supports github for now)
	syncCmd.Flags().
		StringVar(&remote, "remote", "github.com/redds-be/rpkgm-main",
			"Specify a remote, only works with GitHub for now. (ex: github.com/<user>/<repo> without .git)")

	// Flag for the repo's name
	syncCmd.Flags().StringVarP(&name, "name", "n", "main", "Name of the repository.")

	// If you use --remote, you should use --name
	syncCmd.MarkFlagsRequiredTogether("remote", "name", "repo")
}
