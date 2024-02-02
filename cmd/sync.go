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
	"github.com/redds-be/rpkgm/internal/sync"
	"github.com/redds-be/rpkgm/internal/util"
	"github.com/spf13/cobra"
)

var (
	repoName string
	remote   string
)

// syncCmd represents the sync command.
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the main or a specified repo.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if user is root.
		util.CheckRoot()

		// Decide what to do and do what is needed to do
		sync.Decide(repoDB, importFile, remote, repoName)
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
	syncCmd.Flags().StringVarP(&repoName, "name", "n", "main", "Name of the repository.")

	// If you use --remote, you should use --name
	syncCmd.MarkFlagsRequiredTogether("remote", "name", "repo")
}
