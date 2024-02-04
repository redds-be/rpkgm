package sync

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/redds-be/rpkgm/internal/add"
	"github.com/redds-be/rpkgm/internal/database"
	"github.com/redds-be/rpkgm/internal/util"
)

// dlFromRemote downloads the repo's JSON file and the packages build files from a remote.
func dlFromRemote(remote, repoName string) string {
	destDir := fmt.Sprintf("var/rpkgm/%s", repoName)

	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not create the destination directory for the build files.",
		)
		os.Exit(1)
	}

	archive := fmt.Sprintf("var/rpkgm/%s/%s.tar.gz", repoName, repoName)

	err = util.Download(
		archive,
		fmt.Sprintf("https://%s/raw/main/%s.tar.gz", remote, repoName),
	)
	if err != nil {
		util.Display(os.Stderr, true, "rpkgm could not download the build files.")
		os.Exit(1)
	}

	importFile := fmt.Sprintf("var/rpkgm/%s/repo.json", repoName)

	err = util.Download(
		importFile,
		fmt.Sprintf("https://%s/raw/main/repo.json", remote),
	)
	if err != nil {
		util.Display(os.Stderr, true, "rpkgm could not download the JSON file of the repo.")
		os.Exit(1)
	}

	_, err = util.Untar(destDir, archive)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not untar the repo's archive. Error: %s",
			err,
		)
		os.Exit(1)
	}

	return importFile
}

// syncWithFile syncs a repo using a file.
func syncWithFile(importFile string, dbAdapter *database.Adapter) { //nolint:funlen,cyclop
	// Open the file to import
	jsonPkgFile, err := os.Open(importFile)
	if err != nil {
		util.Display(os.Stderr, true, "rpkgm couldn't open the json file. Error: %s", err)
		os.Exit(1)
	}

	// Read the file's content
	contentInbytes, err := io.ReadAll(jsonPkgFile)
	if err != nil {
		util.Display(
			os.Stderr, true,
			"rpkgm couldn't read the json file's content. Error: %s",
			err,
		)
		os.Exit(1)
	}

	// Initialize the Packages struct
	var pkgs database.Packages

	// Read the json content of the file
	err = json.Unmarshal(contentInbytes, &pkgs)
	if err != nil {
		util.Display(os.Stderr, true, "rpkgm couldn't read the json file. Error: %s", err)
		os.Exit(1)
	}

	// for every package in the json file, update their record
	for index := 0; index < len(pkgs.Packages); index++ {
		// Since UPDATE will set null strings as null in the database,
		// we need to get the old information to avoid deleting some information if the JSON field is null.
		pkgInfo, err := dbAdapter.GetPkgInfo(pkgs.Packages[index].Name)
		if err != nil {
			util.Display(
				os.Stderr,
				true,
				"rpkgm was unable to get the information for the package %s to compare against new information. Error: %s",
				pkgs.Packages[index].Name,
				err,
			)

			continue
		}

		// If there isn't a description, give the previous one by default
		if pkgs.Packages[index].Description == "" {
			pkgs.Packages[index].Description = pkgInfo.Description
		}

		// If there isn't a build files dir, give the previous one by default
		if pkgs.Packages[index].BuildFilesDir == "" {
			pkgs.Packages[index].BuildFilesDir = pkgInfo.BuildFilesDir
		}

		// Remove any trailing /
		pkgs.Packages[index].BuildFilesDir = strings.TrimSuffix(
			pkgs.Packages[index].BuildFilesDir,
			"/",
		)

		// If there isn't a repo version, give the previous one by default
		if pkgs.Packages[index].Version == "" {
			pkgs.Packages[index].Version = pkgInfo.RepoVersion
		}

		// If there isn't an archive url, give the previous one by default
		if pkgs.Packages[index].ArchiveURL == "" {
			pkgs.Packages[index].ArchiveURL = pkgInfo.ArchiveURL
		}

		// If there isn't an archive's hash, give the previous one by default
		if pkgs.Packages[index].Sha512 == "" {
			pkgs.Packages[index].Sha512 = pkgInfo.Sha512
		}

		// If there isn't a dependencies list, give the previous one by default
		if pkgs.Packages[index].Dependencies == "" {
			pkgs.Packages[index].Dependencies = pkgInfo.Dependencies
		}

		// Add the package to the repo
		err = dbAdapter.SyncRepo(
			pkgs.Packages[index].Name,
			pkgs.Packages[index].Description,
			pkgs.Packages[index].Version,
			pkgs.Packages[index].BuildFilesDir,
			pkgs.Packages[index].ArchiveURL,
			pkgs.Packages[index].Sha512,
			pkgs.Packages[index].Dependencies,
		)
		if err != nil {
			util.Display(
				os.Stderr, true,
				"rpkgm was unable to update %s in the repo. Error: %s",
				pkgs.Packages[index].Name,
				err,
			)
		}
	}

	// Close the json file
	err = jsonPkgFile.Close()
	if err != nil {
		util.Display(os.Stderr, true, "rpgkm couln't close the json file. Error: %s", err)
		os.Exit(1)
	}
}

// Decide decides what to do based on the given strings.
func Decide(repoDB, importFile, remote, repoName string) {
	if importFile == "" {
		importFile = dlFromRemote(remote, repoName)
	}

	var doAdd bool

	if _, err := os.Stat(repoDB); errors.Is(err, os.ErrNotExist) {
		doAdd = true
	}

	// Connect to the database
	dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not connect to the database. Error: %s",
			err,
		)
		os.Exit(1)
	}

	// Create the packages table if it does not exist
	err = dbAdapter.CreatePkgTable()
	if err != nil {
		util.Display(
			os.Stderr,
			true,
			"rpkgm could not create the packages table in the repo. Error: %s",
			err,
		)
		os.Exit(1)
	}

	if doAdd {
		add.ImportPkgs(importFile, dbAdapter)
	} else {
		syncWithFile(importFile, dbAdapter)
	}

	// Close the database connection
	err = dbAdapter.CloseDBConnection()
	if err != nil {
		util.Display(
			os.Stderr, true,
			"rpkgm could not close the connection to the database. Error: %s",
			err,
		)
		os.Exit(1)
	}
}
