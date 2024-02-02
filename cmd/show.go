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

	"github.com/redds-be/rpkgm/internal/show"
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
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show license related informations and packages general informations.",
	Run: func(cmd *cobra.Command, args []string) {
		// show rpkgm's warranty notice
		if showNotice {
			util.Display(os.Stdout, false, "%s", notice)
			os.Exit(0)
		}

		if showLicense || showInfo {
			if name == "" {
				util.Display(
					os.Stderr,
					false,
					"You need to specify a package name with --name/-n before showing its information.",
				)
				err := cmd.Help()
				if err != nil {
					util.Display(
						os.Stderr,
						false,
						"rpkgm could not display the help message. Error: %s",
						err,
					)
				}
				os.Exit(1)
			}
		}

		// Decide what to do and do what is needed to do
		show.Decide(repoDB, name, showLicense, showInfo, showAll)
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
	showCmd.Flags().
		BoolVarP(&showInfo, "info", "i", false, "Show a given package's general information.")

	// Flag to show every packages general information.
	showCmd.Flags().
		BoolVarP(&showAll, "all", "a", false, "Show every package's general information.")

	// Optional flag to specify repo database location
	showCmd.Flags().
		StringVarP(&repoDB, "repo", "r", "var/rpkgm/main/main.db", "Specify repo Database location.")
}
