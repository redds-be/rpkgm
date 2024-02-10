package update

import (
	"os"

	"github.com/redds-be/rpkgm/internal/database"
	"github.com/redds-be/rpkgm/internal/pkg"
	"github.com/redds-be/rpkgm/internal/util"
)

// checkUpdate checks every installed package to see if there's an update (only informative).
func checkUpdate(installPkgsInfo []database.PkgInfo) {
	var isThereAnUpdate bool
	// For each installed, package, check if the repo's version if different and inform the user
	for _, pkgInfo := range installPkgsInfo {
		if pkgInfo.InstalledVersion != pkgInfo.RepoVersion {
			util.Display(
				os.Stdout,
				false,
				"Update available for %s, Current version: %s | New version: %s",
				pkgInfo.Name,
				pkgInfo.InstalledVersion,
				pkgInfo.RepoVersion,
			)
			isThereAnUpdate = true
		}
	}

	// If there isn't any update, also inform the user
	if !isThereAnUpdate {
		util.Display(os.Stdout, false, "No new updates available.")
	}
}

// Decide decides what to do based on the booleans.
func Decide( //nolint:funlen,gocognit,cyclop
	repoDB string,
	packageList []string,
	all, check, verbose, yes, keep bool,
) {
	// Connect to the database
	dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
	if err != nil {
		util.Display(os.Stderr, true, "rpkgm could not connect to the database. Error: %s", err)
		os.Exit(1)
	}

	// If we check or we want to update every packages, get their info
	var installedPkgsInfo []database.PkgInfo
	if check || all {
		// Get every installed packages's general information
		installedPkgsInfo, err = dbAdapter.GetInstalledPkgInfo()
		if err != nil {
			util.Display(
				os.Stderr,
				true,
				"rpgkm get every installed packages's general information. Error: %s",
				err,
			)
		}
	}

	// If we update every package, append packages that have an update to packageList
	if all {
		for _, pkgInfo := range installedPkgsInfo {
			if pkgInfo.InstalledVersion != pkgInfo.RepoVersion {
				packageList = append(packageList, pkgInfo.Name)
			}
		}
	}

	// If packageList isn't empty, check if they are installed, get their info,
	// ask before updating, and update them
	if len(packageList) > 0 { //nolint:nestif
		for index, pkgName := range packageList {
			// Check if the package is installed
			isInstalled, err := dbAdapter.IsInstalled(pkgName)
			if err != nil {
				util.Display(
					os.Stderr,
					true,
					"rpkgm could not determine if %s is installed or not. Error: %s",
					pkgName,
					err,
				)
				os.Exit(1)
			}

			// If it isn't installed, skip
			if !isInstalled {
				util.Display(
					os.Stderr,
					true,
					"The package named %s is not installed. Skipping...",
					pkgName,
				)

				continue
			}

			// Get the package's general info
			pkgInfo, err := dbAdapter.GetPkgInfo(pkgName)
			if err != nil {
				util.Display(
					os.Stderr,
					true,
					"rpkgm could not get %s's general information. Error: %s",
					pkgName,
					err,
				)
				os.Exit(1)
			}

			// If its installed and if there's an update, update it
			if isInstalled && pkgInfo.InstalledVersion != pkgInfo.RepoVersion {
				util.Display(
					os.Stdout,
					true,
					"Updating %s from version %s to version %s.",
					pkgInfo.Name,
					pkgInfo.InstalledVersion,
					pkgInfo.RepoVersion,
				)

				// Append to the marked packages
				pkg.MarkedPkgs = append(pkg.MarkedPkgs, pkgName)

				// Ask
				if !yes {
					pkg.Ask(dbAdapter)
				}

				// Install
				err = pkg.Install(
					pkgInfo,
					index+1,
					len(packageList),
					verbose,
					keep,
					false,
					dbAdapter,
				)
				if err != nil {
					util.Display(os.Stderr, true, "%s", err)

					continue
				}
			} else {
				util.Display(os.Stdout, true, "No updates available for %s.", pkgName)
			}
		}
	}

	// If we check for updates, call checkUpdate
	if check {
		if len(installedPkgsInfo) > 0 {
			checkUpdate(installedPkgsInfo)
		}
	}

	// Close the connection to the database
	err = dbAdapter.CloseDBConnection()
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not close the connection to the database. Error: %s",
			err,
		)
	}
}
