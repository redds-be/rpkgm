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

	"github.com/redds-be/rpkgm/internal/pkg"
	"github.com/redds-be/rpkgm/internal/util"
	"github.com/spf13/cobra"
)

var (
	toInstall   []string
	toUninstall []string
	verbose     bool
	keep        bool
	force       bool
	yes         bool
	repoDB      string
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
		if len(toInstall) > 0 || len(toUninstall) > 0 {
			// Send everything to pkg.Manage (not a fan of this, will change later)
			pkg.Manage(toInstall, toUninstall, verbose, keep, force, yes, repoDB)
		} else {
			err := cmd.Help()
			if err != nil {
				util.Display(os.Stderr, "rpkgm could not display the help message. Error: %s", err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		util.Display(os.Stderr, "rpkgm could not start. Error: %s", err)
	}
}

// init initializes the command-line arguments for cobra.
func init() { //nolint:gochecknoinits
	// Flag for a list of packages to install
	rootCmd.Flags().
		StringSliceVarP(&toInstall, "install", "i", nil, "Package(s) to install. For multiple packages, separate them with commas.")

	// Flag for a list of package to remove
	rootCmd.Flags().
		StringSliceVarP(&toUninstall, "uninstall", "u", nil, "Package(s) to uninstall. For multiple packages, separate them with commas.")

	// Mark install and remove as mutually exclusive
	rootCmd.MarkFlagsMutuallyExclusive("install", "uninstall")

	// Flag for verbosity
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Make rpkgm verbose during operation.")

	// Flag for keeping packages source dir intact after installation
	rootCmd.Flags().
		BoolVarP(&keep, "keep", "k", false, "Keep package(s) source directories after installation (/usr/src/rpkgm/<pkgName>)")

	// Flag for forcing the installation of already installed packages
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "Force install already installed packages.")

	// Flag to indicate there is no need for confirmation
	rootCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Do not ask before installing/uninstalling.")

	// Optional flag to specify repo database location
	rootCmd.Flags().
		StringVarP(&repoDB, "repo", "r", "var/rpkgm/main.db", "Specify repo Database location (Defaults to /var/rpkgm/main.db)")
}
