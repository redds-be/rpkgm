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

package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"

	"github.com/redds-be/rpkgm/internal/logging"
)

// Define some colors.
var (
	Rc = "\033[0m"
	By = "\033[1m\033[33m"
	Bg = "\033[1m\033[32m"
)

// CheckRoot checks if the user is root.
func CheckRoot() {
	currUser, err := user.Current()
	if err != nil {
		Display(os.Stderr, true, "Unable to determine if rpkgm is running as root.")
		os.Exit(1)
	}

	if currUser.Uid != "0" {
		Display(os.Stderr, false, "Please run rpkgm as root.")
		os.Exit(1)
	}
}

// Display is wrapper over fmt.Fprintf.
func Display(out io.Writer, doLog bool, format string, toDisplay ...any) {
	_, err := fmt.Fprintf(out, fmt.Sprintf("%s\n", format), toDisplay...)
	if err != nil {
		log.Println("rpkgm was unable to print output...")
	}

	if doLog {
		logging.LogToFile(format, toDisplay...)
	}
}
