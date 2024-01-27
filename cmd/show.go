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

	"github.com/spf13/cobra"
)

var showNotice bool //nolint:gochecknoglobals

const notice = `rpgkm: redd's package manager.
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
var showCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "show",
	Short: "Show license related informations",
	Run: func(cmd *cobra.Command, args []string) {
		if showNotice {
			fmt.Println(notice) //nolint:forbidigo
		} else {
			err := cmd.Help()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

// init defines flags.
//
//nolint:gochecknoinits
func init() {
	// Link to root (root = 'rpkgm', show = 'rpkgm show')
	rootCmd.AddCommand(showCmd)

	// --warranty or -w is used to print the GPLv3 warranty notice.
	showCmd.Flags().BoolVarP(&showNotice, "warranty", "w", false, "Show warranty.")
}
