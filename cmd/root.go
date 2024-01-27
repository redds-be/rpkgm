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
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	toInstall []string //nolint:gochecknoglobals
	toRemove  []string //nolint:gochecknoglobals
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "rpkgm",
	Short: "rpkgm Copyright (C) 2024 redd.",
	Long: `rpkgm Copyright (C) 2024 redd
This program comes with ABSOLUTELY NO WARRANTY; for details type 'rpkgm show -w'.
This is free software, and you are welcome to redistribute it
under certain conditions; see <https://www.gnu.org/licenses/gpl-3.0.html>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(toInstall) > 0 {
			for _, pkg := range toInstall {
				log.Println("Installing", pkg)
			}
		} else if len(toRemove) > 0 {
			for _, pkg := range toRemove {
				log.Println("Removing", pkg)
			}
		} else {
			err := cmd.Help()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init defines flags.
//
//nolint:gochecknoinits
func init() {
	// Flag for a list of packages to install
	rootCmd.Flags().
		StringSliceVarP(&toInstall, "install", "i", nil, "Package(s) to install. For multiple packages, separate them with commas.")

	// Flag for a list of package to remove
	rootCmd.Flags().
		StringSliceVarP(&toRemove, "remove", "r", nil, "Package(s) to remove. For multiple packages, separate them with commas.")

	// Mark install and remove as mutually exclusive
	rootCmd.MarkFlagsMutuallyExclusive("install", "remove")
}
