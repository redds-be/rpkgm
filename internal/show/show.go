package show

import (
	"fmt"
	"io"
	"os"

	"github.com/redds-be/rpkgm/internal/database"
	"github.com/redds-be/rpkgm/internal/util"
)

// printLicense prints the license of a given package assuming the license if in the build files.
func printLicense(name string, dbAdapter *database.Adapter) {
	// Find the build files for the given package, the license should be in there
	buildFilesDir, err := dbAdapter.GetPkgBuildFilesDir(name)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not find the build files (build files includes the license) for the given package. Error: %s",
			err,
		)
		os.Exit(1)
	}

	// Open the license file
	licenseFile, err := os.Open(fmt.Sprintf("%s/LICENSE", buildFilesDir))
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not open or find the license file for the given package. Error: %s",
			err,
		)
		os.Exit(1)
	}

	// Read the license file
	licenseContent, err := io.ReadAll(licenseFile)
	if err != nil {
		util.Display(
			os.Stderr, true,
			"rpkgm could not read from the license file of the given package. Error: %s",
			err,
		)
		os.Exit(1)
	}

	// Print the license file's content
	util.Display(os.Stdout, false, "%v", string(licenseContent))

	// Close the license file
	err = licenseFile.Close()
	if err != nil {
		util.Display(
			os.Stderr, true,
			"rpkgm could not close the license file for the given package. Error: %s",
			err,
		)
		os.Exit(1)
	}
}

// printInfo prints the general information of a given package.
func printInfo(name string, dbAdapter *database.Adapter) {
	// Get the given package's general info
	pkgInfo, err := dbAdapter.GetPkgInfo(name)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not get the given package's information. Error: %s",
			err,
		)
		os.Exit(1)
	}

	// Display the general info (changes depending on the installation status)
	if pkgInfo.Installed {
		util.Display(
			os.Stdout, false,
			"%s [Installed (%s)]\t- %s\t- Repo's version: %s",
			pkgInfo.Name,
			pkgInfo.InstalledVersion,
			pkgInfo.Description,
			pkgInfo.RepoVersion,
		)
	} else {
		util.Display(os.Stdout, false, "%s [Not installed]\t- %s\t- Repo's version: %s", pkgInfo.Name, pkgInfo.Description, pkgInfo.RepoVersion)
	}
}

// printAllinfo prints the general information of every package in the repo.
func printAllinfo(dbAdapter *database.Adapter) {
	// Get every package info in the database
	allPkgInfo, err := dbAdapter.GetAllPkgInfo()
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could query the repo's database. Error: %s",
			err,
		)
		os.Exit(1)
	}

	// For every package, display the info (changes depending on the installation status)
	for _, pkgInfo := range allPkgInfo {
		if pkgInfo.Installed {
			util.Display(
				os.Stdout, false,
				"%s [Installed (%s)]\t- %s\t- Repo's version: %s",
				pkgInfo.Name,
				pkgInfo.InstalledVersion,
				pkgInfo.Description,
				pkgInfo.RepoVersion,
			)
		} else {
			util.Display(os.Stdout, false, "%s [Not installed]\t- %s\t- Repo's version: %s", pkgInfo.Name, pkgInfo.Description, pkgInfo.RepoVersion)
		}
	}
}

// Decide decies what to do based on the given booleans.
func Decide(repoDB, name string, showLicense, showInfo, showAll bool) {
	// Connect to the database
	dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could connect to the repo's database. Error: %s",
			err,
		)
		os.Exit(1)
	}

	// Show the license of the package
	if showLicense {
		printLicense(name, dbAdapter)
	}

	// Show the general info of the package
	if showInfo {
		printInfo(name, dbAdapter)
	}

	// Show the general info of all the packages
	if showAll {
		printAllinfo(dbAdapter)
	}

	// Close the database connection
	err = dbAdapter.CloseDBConnection()
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not close the connection to the database. Error: %s",
			err,
		)
		os.Exit(1)
	}
}
