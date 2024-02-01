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
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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

// Download downloads a body from a url and writes to dest.
func Download(dest, url string) error {
	// Create the destination file
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}

	// Prepare the client and the request
	client := &http.Client{}
	ctx := context.Background()
	dlReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	// Do the request
	resp, err := client.Do(dlReq)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("url returned code %d", resp.StatusCode) //nolint:goerr113
	}

	// Copy the content of the body to the newly created file
	_, err = io.Copy(destFile, resp.Body)
	if err != nil {
		return err
	}

	// Close the destination file
	err = destFile.Close()
	if err != nil {
		return err
	}

	// Close the body
	err = resp.Body.Close()

	return err
}

func Untar(destDir, archive string) error { //nolint:cyclop
	tarFile, err := os.Open(archive)
	if err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(tarFile)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(fmt.Sprintf("%s/%s", destDir, header.Name), os.ModePerm); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(fmt.Sprintf("%s/%s", destDir, header.Name))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil { //nolint:gosec
				return err
			}

			err = outFile.Close()
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf( //nolint:goerr113
				"unknown type: %v in %v",
				header.Typeflag,
				header.Name,
			)
		}
	}

	err = tarFile.Close()
	if err != nil {
		return err
	}

	err = gzipReader.Close()

	return err
}
