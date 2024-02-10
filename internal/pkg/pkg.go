package pkg

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/redds-be/rpkgm/internal/database"
	"github.com/redds-be/rpkgm/internal/util"
)

// MarkedPkgs is a slice that contains the name of the packages marked for an operation.
var MarkedPkgs []string

// resolveDeps recursively resolves the dependencies of a slice of dependencies and marks them for installation.
func resolveDeps(mainPkgName string, deps []string, dbAdapter *database.Adapter) { //nolint:funlen
	for _, pkgName := range deps {
		// Since sometimes the deps list contains a single empty string (due to using split on an empty string)
		// we just ignore it
		if pkgName == "" {
			continue
		}

		// Check if the dependency is in the repo
		isInRepo, _ := dbAdapter.IsPkgInRepo(pkgName)
		if !isInRepo {
			util.Display(
				os.Stderr,
				true,
				"%s's dependency named %s is not in the repository, you will need to install it yourself...",
				mainPkgName,
				pkgName,
			)

			continue
		}

		// Get the information for the package
		pkgInfo, err := dbAdapter.GetPkgInfo(pkgName)
		if err != nil {
			util.Display(
				os.Stderr,
				true,
				"rpkgm couldn't get %s's dependency %s's information despite it being in the repo, you will need to install it yourself...",
				mainPkgName,
				pkgName,
			)

			continue
		}

		// Check if the dependency is already installed
		isInstalled, err := dbAdapter.IsInstalled(pkgName)
		if err != nil {
			util.Display(
				os.Stderr,
				true,
				"rpkgm couldn't determine if %s's dependency %s is already installed, you will need to install it yourself...",
				mainPkgName,
				pkgName,
			)

			continue
		}

		// Check if the dependency is a duplicate
		isDuplicate := slices.Contains(MarkedPkgs, pkgName)

		// If the dependency is in the repo, is not already installed and is not a duplicate,
		// resolve its dependencies and mark de dependency for installation
		if isInRepo && !isInstalled && !isDuplicate {
			moreDeps := strings.Split(pkgInfo.Dependencies, " ")
			if len(moreDeps) > 0 {
				resolveDeps(pkgName, moreDeps, dbAdapter)
			}
			MarkedPkgs = append(MarkedPkgs, pkgName)
			util.Display(os.Stdout, true, "Installing %s=%s", pkgName, pkgInfo.RepoVersion)
		}
	}
}

// Ask asks before doing any operation.
func Ask(dbAdapter *database.Adapter) {
	// If there are marked packages, ask, else, just quit
	if len(MarkedPkgs) > 0 { //nolint:nestif
		var choice string

		// Ask the user's confirmation
		fmt.Printf("Do you want to perform this operation? [y/N] ") //nolint:forbidigo
		_, err := fmt.Scanln(&choice)
		if err != nil {
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
			os.Exit(1)
		}

		// Choice is yes, we go on
		if choice == "y" || choice == "Y" {
			return
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

	// If we're here, there isn't any packages marked for an operation
	util.Display(os.Stderr, true, "No package selected for any operations.")

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
	os.Exit(0)
}

// Install installs a package.
func Install( //nolint:funlen,cyclop
	pkgInfo database.PkgInfo,
	index, total int,
	verbose, keep, force bool,
	dbAdapter *database.Adapter,
) error {
	// Set the destination directory
	destDir := fmt.Sprintf("/tmp/rpkgm/%s", pkgInfo.Name)
	if keep {
		destDir = fmt.Sprintf("/tmp/usr/src/rpkgm/%s", pkgInfo.Name)
	}

	// Remove the destination directory if it already exists
	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		err := os.RemoveAll(destDir)
		if err != nil {
			return fmt.Errorf(
				"rpkgm was unable to pre-clean the build directory of %s, Error: %w",
				pkgInfo.Name,
				err,
			)
		}
	}

	// Create the destination directory
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf(
			"rpkgm was unable to create the build dir for %s, Error: %w",
			pkgInfo.Name,
			err,
		)
	}

	// Set the destination file of the archive
	archive := fmt.Sprintf("%s/%s.tar.gz", destDir, pkgInfo.Name)

	// Inform of the downloading
	util.Display(
		os.Stdout, false,
		"Downloading (%s%d%s of %s%d%s) %s%s=%s%s",
		util.By,
		index,
		util.Rc,
		util.By,
		total,
		util.Rc,
		util.Bg,
		pkgInfo.Name,
		pkgInfo.RepoVersion,
		util.Rc,
	)

	// Download the archive
	err = util.Download(archive, pkgInfo.ArchiveURL)
	if err != nil {
		return fmt.Errorf(
			"rpkgm was unable to download the archive for %s, Error: %w",
			pkgInfo.Name,
			err,
		)
	}

	// If we don't force the installation, verify the downloaded archive
	if !force {
		isOk, err := util.Verify(archive, pkgInfo.Sha512)
		if err != nil {
			return fmt.Errorf(
				"rpkgm was unable to verify the archive for %s, Error: %w",
				pkgInfo.Name,
				err,
			)
		}

		if !isOk {
			return errors.New( //nolint:goerr113
				"the archive's hash does not correspond to the hash in the repo, you can ignore this error by re-running with --force/-f",
			)
		}
	}

	// Inform of the extracting
	util.Display(
		os.Stdout, false,
		"Extracting (%s%d%s of %s%d%s) %s%s=%s%s",
		util.By,
		index,
		util.Rc,
		util.By,
		total,
		util.Rc,
		util.Bg,
		pkgInfo.Name,
		pkgInfo.RepoVersion,
		util.Rc,
	)

	// Untar the archive
	newDestDir, err := util.Untar(destDir, archive)
	if err != nil {
		return fmt.Errorf(
			"rpkgm was unable to extract the archive of %s, Error: %w",
			pkgInfo.Name,
			err,
		)
	}

	// Set the source and destination makefiles
	makefileSrc := fmt.Sprintf("%s/Makefile", pkgInfo.BuildFilesDir)
	makeFileDst := fmt.Sprintf("%s/Makefile", newDestDir)

	// Copy the source make into the destination makefile
	err = util.Copy(makefileSrc, makeFileDst, true)
	if err != nil {
		return fmt.Errorf(
			"rpkgm was unable to copy its own Makefile for %s, Error: %w",
			pkgInfo.Name,
			err,
		)
	}

	// Inform of the installing
	util.Display(
		os.Stdout, false,
		"Installing (%s%d%s of %s%d%s) %s%s=%s%s",
		util.By,
		index,
		util.Rc,
		util.By,
		total,
		util.Rc,
		util.Bg,
		pkgInfo.Name,
		pkgInfo.RepoVersion,
		util.Rc,
	)

	// Install the package
	install := fmt.Sprintf("cd %s && make install", newDestDir)
	inOut, err := exec.Command("/usr/bin/env", "bash", "-c", install).CombinedOutput()
	if err != nil {
		// In case of errors, be verbose to leave a trace
		util.Display(io.Discard, true, "%s", string(inOut))
		if verbose && string(inOut) != "" {
			util.Display(os.Stdout, false, string(inOut))
		}

		return fmt.Errorf(
			"rpkgm was unable to install the package %s, Error: %w",
			pkgInfo.Name,
			err,
		)
	}

	// Little hack that may be removed later to log the output
	util.Display(io.Discard, true, "%s", string(inOut))

	// Display the output
	if verbose && string(inOut) != "" {
		util.Display(os.Stdout, false, string(inOut))
	}

	// If we don't keep the build dir, remove it
	if !keep {
		// Inform of the cleaning
		util.Display(
			os.Stdout, false,
			"Cleaning (%s%d%s of %s%d%s) %s%s=%s%s",
			util.By,
			index,
			util.Rc,
			util.By,
			total,
			util.Rc,
			util.Bg,
			pkgInfo.Name,
			pkgInfo.RepoVersion,
			util.Rc,
		)

		err = os.RemoveAll(destDir)
		if err != nil {
			return fmt.Errorf(
				"rpkgm was unable to clean the build directory of %s, Error: %w",
				pkgInfo.Name,
				err,
			)
		}
	}

	// Mark the installed package as installed
	err = dbAdapter.MarkAsInstalled(pkgInfo.Name)
	if err != nil {
		return fmt.Errorf(
			"rpkgm could not mark %s as installed although the package is, in fact installed. Error: %w",
			pkgInfo.Name,
			err,
		)
	}

	// Set the package's installed version to the repo's version
	err = dbAdapter.SetInstalledVersion(pkgInfo.Name, pkgInfo.RepoVersion)
	if err != nil {
		return fmt.Errorf(
			"rpkgm could not set %s's install version in the repo, although the package is, in fact installed. Error: %w",
			pkgInfo.Name,
			err,
		)
	}

	return nil
}

// uninstall uninstalls a package.
func uninstall( //nolint:funlen
	pkgInfo database.PkgInfo,
	index, total int,
	verbose, keep bool,
	dbAdapter *database.Adapter,
) error {
	// Inform of the uninstalling
	util.Display(
		os.Stdout, false,
		"Uninstalling (%s%d%s of %s%d%s) %s%s=%s%s",
		util.By,
		index,
		util.Rc,
		util.By,
		total,
		util.Rc,
		util.Bg,
		pkgInfo.Name,
		pkgInfo.InstalledVersion,
		util.Rc,
	)

	// Uninstall the package
	uninstall := fmt.Sprintf("cd %s && make uninstall", pkgInfo.BuildFilesDir)
	unOut, err := exec.Command("/usr/bin/env", "bash", "-c", uninstall).CombinedOutput()
	if err != nil {
		// In case of errors, be verbose to leave a trace
		util.Display(io.Discard, true, "%s", string(unOut))
		if verbose && string(unOut) != "" {
			util.Display(os.Stdout, false, string(unOut))
		}

		return fmt.Errorf(
			"rpkgm was unable to uninstall the package %s, Error: %w",
			pkgInfo.Name,
			err,
		)
	}

	// If we don't keep the source, remove it
	if !keep {
		workdir := fmt.Sprintf("/tmp/usr/src/rpkgm/%s", pkgInfo.Name)
		if _, err := os.Stat(workdir); !os.IsNotExist(err) {
			// Inform of the cleaning
			util.Display(
				os.Stdout, false,
				"Cleaning (%s%d%s of %s%d%s) %s%s=%s%s",
				util.By,
				index,
				util.Rc,
				util.By,
				total,
				util.Rc,
				util.Bg,
				pkgInfo.Name,
				pkgInfo.InstalledVersion,
				util.Rc,
			)

			// Clean working directory
			err = os.RemoveAll(workdir)
			if err != nil {
				return fmt.Errorf(
					"rpkgm could not remove %s's build directory. Error: %w",
					pkgInfo.Name,
					err,
				)
			}
		}
	}

	// Mark the package as not installed anymore
	err = dbAdapter.MarkAsNotInstalled(pkgInfo.Name)
	if err != nil {
		return fmt.Errorf(
			"rpkgm could not mark %s as not installed although the package is, in fact uninstalled. Error: %w",
			pkgInfo.Name,
			err,
		)
	}

	// Remove the package installed version in the repo
	err = dbAdapter.SetInstalledVersion(pkgInfo.Name, "")
	if err != nil {
		return fmt.Errorf(
			"rpkgm could not remove %s's installed version (in the repo's db) although the package is, in fact uninstalled. Error: %w",
			pkgInfo.Name,
			err,
		)
	}

	return nil
}

// Decide decides what to do based on the given booleans.
func Decide( //nolint:funlen,gocognit,cyclop
	doInstall, force, verbose, keep, yes, resolve bool,
	packageList []string,
	repoDB string,
) {
	// Check if the user is root
	util.CheckRoot("Please run rpkgm as root.")

	// Connect to the database
	dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
	if err != nil {
		util.Display(os.Stderr, false, "rpkgm could not connect to the database. Error: %s", err)
		os.Exit(1)
	}

	for _, pkgName := range packageList {
		// Check if the package is in the repo
		isInRepo, _ := dbAdapter.IsPkgInRepo(pkgName)
		if !isInRepo && !force {
			util.Display(
				os.Stderr,
				true,
				"The package named %s is not in the repository, skipping...",
				pkgName,
			)

			continue
		}

		// Get the package's general information
		pkgInfo, err := dbAdapter.GetPkgInfo(pkgName)
		if err != nil {
			util.Display(
				os.Stderr,
				true,
				"rpkgm couldn't get %s's information despite it being in the repo, skipping...",
				pkgName,
			)

			continue
		}

		// Check if the package is already installed
		isInstalled, err := dbAdapter.IsInstalled(pkgName)
		if err != nil {
			util.Display(
				os.Stderr,
				true,
				"rpkgm couldn't determine if %s is installed or not, skipping...",
				pkgName,
			)

			continue
		}

		// Get the dependencies as a list
		deps := strings.Split(pkgInfo.Dependencies, " ")

		switch {
		// Case the operation is installing, the package is already installed but we don't force the re-installation it, skip it
		case doInstall && isInstalled && !force:
			util.Display(
				os.Stdout,
				true,
				"The package named %s is already installed, you can force re-installation by re-running using --force/-f. Skipping...",
				pkgName,
			)

			continue
			// Case the operation is installing, the package is already installed and we force the re-installation, we install it
		case doInstall && isInstalled && force:
			// If the experimental resolve feature is set, resolve its deps
			if resolve {
				resolveDeps(pkgName, deps, dbAdapter)
			} else {
				util.Display(
					os.Stdout,
					true,
					"The package %s lists %v as dependencies, resolving dependencies is still experimental, either use --resolve or install them yourself.",
					pkgName,
					deps,
				)
			}
			util.Display(os.Stdout, true, "Installing %s=%s", pkgName, pkgInfo.RepoVersion)
			// Mark the package for installation
			MarkedPkgs = append(MarkedPkgs, pkgName)
			// Case the operation is installing and the package is not installed
		case doInstall && !isInstalled:
			// If the experimental resolve feature is set, resolve its deps
			if resolve {
				resolveDeps(pkgName, deps, dbAdapter)
			} else {
				util.Display(
					os.Stdout,
					true,
					"The package %s lists %v as dependencies, resolving dependencies is still experimental, either use --resolve or install them yourself.",
					pkgName,
					deps,
				)
			}
			util.Display(os.Stdout, true, "Installing %s=%s", pkgName, pkgInfo.RepoVersion)
			// Mark the package for installation
			MarkedPkgs = append(MarkedPkgs, pkgName)
			// Case the operation is uninstallation and the package is not installed, we skip it
		case !doInstall && !isInstalled:
			util.Display(
				os.Stdout,
				true,
				"The package named %s is not installed, Skipping...",
				pkgName,
			)

			continue
			// Case the operation is uninstallation and the package is installed, we mark it for uninstallation
		case !doInstall && isInstalled:
			util.Display(os.Stdout, true, "Uninstalling %s=%s", pkgName, pkgInfo.RepoVersion)
			MarkedPkgs = append(MarkedPkgs, pkgName)
			// Case I don't know what the f to do based on the given information
		default:
			util.Display(
				os.Stdout,
				true,
				"rpkgm doesn't know what to do for %s, Skipping...",
				pkgName,
			)

			continue
		}
	}

	// If --yes/-y is not set, we ask before doing anything
	if !yes {
		Ask(dbAdapter)
	}

	for index, pkgName := range MarkedPkgs {
		// Get the general information of the marked package
		pkgInfo, err := dbAdapter.GetPkgInfo(pkgName)
		if err != nil {
			util.Display(
				os.Stderr,
				true,
				"rpkgm couldn't get %s's information despite it being in the repo, skipping...",
				pkgName,
			)

			continue
		}

		// If the operation is installation, call install
		if doInstall {
			err = Install(pkgInfo, index+1, len(MarkedPkgs), verbose, keep, force, dbAdapter)
			if err != nil {
				util.Display(os.Stderr, true, "%s", err)
			}
		}

		// If the operation is uninstallation, call uninstall
		if !doInstall {
			err = uninstall(pkgInfo, index+1, len(MarkedPkgs), verbose, keep, dbAdapter)
			if err != nil {
				util.Display(os.Stderr, true, "%s", err)
			}
		}
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
