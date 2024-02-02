package manage

import (
	"os"

	"github.com/redds-be/rpkgm/internal/database"
	"github.com/redds-be/rpkgm/internal/util"
)

// remove removes a package.
func remove(name string, dbAdapter *database.Adapter) {
	err := dbAdapter.RemovePackage(name)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not delete the package from the repository. Error: %s",
			err,
		)
		os.Exit(1)
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

	os.Exit(0)
}

// changeDesc changes the description of a package.
func changeDesc(name, newDesc string, dbAdapter *database.Adapter) {
	err := dbAdapter.ChangePkgDesc(name, newDesc)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not change the package's description in the repo's database. Error: %s",
			err,
		)
		os.Exit(1)
	}
}

// markAsInstalled marks a package as installed.
func markAsInstalled(name string, dbAdapter *database.Adapter) {
	err := dbAdapter.MarkAsInstalled(name)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not mark the package as installed in the repo's database. Error: %s",
			err,
		)
		os.Exit(1)
	}
}

// markAsNotInstalled marks a package as not installed.
func markAsNotInstalled(name string, dbAdapter *database.Adapter) {
	err := dbAdapter.MarkAsNotInstalled(name)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not mark the package as not installed in the repo's database. Error: %s",
			err,
		)
		os.Exit(1)
	}
}

// changeInstalledVersion changes the installed version of a package.
func changeInstalledVersion(name, installedVersion string, dbAdapter *database.Adapter) {
	err := dbAdapter.SetInstalledVersion(name, installedVersion)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not set the package's installed version. Error: %s",
			err,
		)
		os.Exit(1)
	}
}

// changeRepoVersion changes the repo version of package.
func changeRepoVersion(name, repoVersion string, dbAdapter *database.Adapter) {
	err := dbAdapter.UpdateRepoVersion(name, repoVersion)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not set the package's repo version. Error: %s",
			err,
		)
		os.Exit(1)
	}
}

// rename renames a package.
func rename(name, newName string, dbAdapter *database.Adapter) {
	err := dbAdapter.RenamePackage(name, newName)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not rename the package in the repo's database. Error: %s",
			err,
		)
		os.Exit(1)
	}
}

// Decide decies what to do based on the given booleans.
func Decide( //nolint:funlen,cyclop
	repoDB, name, newName, newDesc, installedVersion, repoVersion string,
	doRemove, markInstalled, markNotInstalled bool,
) {
	// Connect to the database
	dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
	if err != nil {
		util.Display(os.Stderr, true, "rpkgm could not connect to the database. Error: %s", err)
		os.Exit(1)
	}

	// Check if the given package is in the repo (forcing the close of the db connection since it's not a fatal error)
	isInRepo, _ := dbAdapter.IsPkgInRepo(name)
	if !isInRepo {
		util.Display(os.Stderr, true, "The package: %s is not in the repository.", name)

		// Close the database connection
		err := dbAdapter.CloseDBConnection()
		if err != nil {
			util.Display(
				os.Stderr,
				true,
				"rpkgm could not close the connection to the database. Error: %s",
				err,
			)
			os.Exit(1)
		}
		os.Exit(1)
	}

	// Remove the package
	if doRemove {
		remove(name, dbAdapter)
	}

	// Change the package's description
	if newDesc != "" {
		changeDesc(name, newDesc, dbAdapter)
	}

	// Mark the package as installed or as not installed
	if markInstalled {
		markAsInstalled(name, dbAdapter)
	} else if markNotInstalled {
		markAsNotInstalled(name, dbAdapter)
	}

	// Change or set the package's installed version
	if installedVersion != "" {
		changeInstalledVersion(name, installedVersion, dbAdapter)
	}

	// Change the package's repo version
	if repoVersion != "" {
		changeRepoVersion(name, repoVersion, dbAdapter)
	}

	// Rename the package
	if newName != "" {
		rename(name, newName, dbAdapter)
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
