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
	"github.com/redds-be/rpkgm/internal/update"
	"github.com/redds-be/rpkgm/internal/util"
	"github.com/spf13/cobra"
)

var (
	updateList []string
	all        bool
)

// updateCmd represents the update command.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "See and install updates.",
	Run: func(cmd *cobra.Command, args []string) {
		// if len(updateList) > 0 = update -u pkg1,pkg2 = update them
		// else if all = update -a = update all
		// else = update = only check updates
		if len(updateList) > 0 {
			// Check if the user is root
			util.CheckRoot("Please run rpkgm update as root.")

			update.Decide(repoDB, updateList, all, false, verbose, yes, keep)
		} else if all {
			// Check if the user is root
			util.CheckRoot("Please run rpkgm update as root.")

			update.Decide(repoDB, nil, true, false, verbose, yes, keep)
		} else {
			update.Decide(repoDB, nil, false, true, false, false, false)
		}
	},
}

func init() { //nolint:gochecknoinits
	// Link to root (root = 'rpkgm', update = 'rpkgm update')
	rootCmd.AddCommand(updateCmd)

	// Flag to update specific packages
	updateCmd.Flags().
		StringSliceVarP(&updateList, "update", "u", nil, "List of package to update separated by commas.")

	// Flag to update every packages that has an update
	updateCmd.Flags().BoolVarP(&all, "all", "a", false, "Update every packages that has an update.")

	// List of packages and every packages are incompatible
	updateCmd.MarkFlagsMutuallyExclusive("update", "all")

	// Flag for verbosity
	updateCmd.Flags().
		BoolVarP(&verbose, "verbose", "v", false, "Make rpkgm verbose during operation.")

	// Flag to indicate there is no need for confirmation
	updateCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Do not ask before updating.")

	// Flag for keeping packages source dir intact after installation
	updateCmd.Flags().
		BoolVarP(&keep, "keep", "k", false, "Keep package(s) source directories after update (/usr/src/rpkgm/<pkgName>)")

	// Optional flag to specify repo database location
	updateCmd.Flags().
		StringVarP(&repoDB, "repo", "r", "var/rpkgm/main/main.db", "Specify repo Database location.")
}
