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
	"crypto/sha512"
	"encoding/hex"
	"errors"
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
		Display(os.Stderr, true, "Unable to determine if rpkgm is running as root. Error: %s", err)
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
		log.Printf("rpkgm was unable to print to output. Error: %s\n", err)
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

// Verify verifies the sha512 hash of a file.
func Verify(fileToVerify, supposedHash string) (bool, error) {
	// Open the file to verify
	file, err := os.Open(fileToVerify)
	if err != nil {
		return false, err
	}

	// Create the hash of the file
	hash := sha512.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, err
	}

	// Close the file
	if err = file.Close(); err != nil {
		return false, err
	}

	// Convert the hash to a hex string
	hexHash := hex.EncodeToString(hash.Sum(nil))

	// Compare the hashes
	if hexHash == supposedHash {
		return true, nil
	}

	return false, nil
}

// Untar untars a tar.gz tarball.
func Untar(destDir, archive string) (string, error) { //nolint:cyclop,funlen
	archiveParentDir := ""

	// Open the tarball
	tarFile, err := os.Open(archive)
	if err != nil {
		return "", err
	}

	// Create a reader to gunzip the tarball
	gzipReader, err := gzip.NewReader(tarFile)
	if err != nil {
		return "", err
	}

	// Create a reader for the tarball
	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if archiveParentDir == "" {
				archiveParentDir = fmt.Sprintf("%s/%s", destDir, header.Name)
			}
			// If the header indicates a dir, create it
			if err := os.MkdirAll(fmt.Sprintf("%s/%s", destDir, header.Name), os.ModePerm); err != nil {
				return "", err
			}
		case tar.TypeReg:
			// If the header indicates a file, create the destination file
			outFile, err := os.Create(fmt.Sprintf("%s/%s", destDir, header.Name))
			if err != nil {
				return "", err
			}

			// Copy the file from the tarball to the newly created file
			if _, err := io.Copy(outFile, tarReader); err != nil { //nolint:gosec
				return "", err
			}
			if err != nil {
				return "", err
			}

			// Close the file
			err = outFile.Close()
			if err != nil {
				return "", err
			}
		case 103: //nolint:gomnd
			// Just ignore a weird header unique to GitHub release tarball
			continue
		default:
			// If it's neither a file nor a dir, what the f is it?
			return "", fmt.Errorf( //nolint:goerr113
				"unknown type: %v in %v",
				header.Typeflag,
				header.Name,
			)
		}
	}

	// Close the tarball
	err = tarFile.Close()
	if err != nil {
		return "", err
	}

	// Close the gzip reader
	err = gzipReader.Close()

	return archiveParentDir, err
}

// Copy copies a file (src) to a new one (dst).
func Copy(src, dst string, overwrite bool) error {
	// If we want to overwrite
	if overwrite {
		// Check if the dst already exists, remove it if it exists
		if _, err := os.Stat(dst); !errors.Is(err, os.ErrNotExist) {
			err = os.Remove(dst)
			if err != nil {
				return err
			}
		}
	}

	// Open the source file
	source, err := os.Open(src)
	if err != nil {
		return err
	}

	// Create the destination file
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}

	// Copy the source file into the destination file
	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	// Close the source file
	err = source.Close()
	if err != nil {
		return err
	}

	// Close the destination file
	err = destination.Close()

	return err
}
