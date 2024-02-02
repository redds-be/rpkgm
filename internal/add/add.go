package add

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/redds-be/rpkgm/internal/database"
	"github.com/redds-be/rpkgm/internal/util"
)

// importPkgs imports packages from a file to a repo.
func ImportPkgs(importFile string, dbAdapter *database.Adapter) { //nolint:funlen
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

	// for every package in the json file, add it to the repo
	for index := 0; index < len(pkgs.Packages); index++ {
		// If there isn't a description, give one by default
		if pkgs.Packages[index].Description == "" {
			pkgs.Packages[index].Description = "[No description provided for this package.]"
		}

		// If there isn't a build files dir, give one by default
		if pkgs.Packages[index].BuildFilesDir == "" {
			pkgs.Packages[index].BuildFilesDir = fmt.Sprintf(
				"var/rpkgm/main/%s",
				pkgs.Packages[index].Name,
			)
		}

		// Remove any trailing /
		pkgs.Packages[index].BuildFilesDir = strings.TrimSuffix(
			pkgs.Packages[index].BuildFilesDir,
			"/",
		)

		// Add the package to the repo
		err = dbAdapter.AddToRepo(
			pkgs.Packages[index].Name,
			pkgs.Packages[index].Description,
			pkgs.Packages[index].Version,
			pkgs.Packages[index].BuildFilesDir,
		)
		if err != nil {
			util.Display(
				os.Stderr, true,
				"rpkgm was unable to add %s to the repo. Error: %s",
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

// addPkg adds a packages with its general information to a repo.
func addPkg(name, description, version, buildFilesDir string, dbAdapter *database.Adapter) {
	// Default value for buildFilesDir (doing it here instead of Flags() because I need 'name')
	if buildFilesDir == "" {
		buildFilesDir = fmt.Sprintf("var/rpkgm/main/%s", name)
	}

	// Remove any trailing /
	buildFilesDir = strings.TrimSuffix(buildFilesDir, "/")

	// Add the package to the main repo
	err := dbAdapter.AddToRepo(name, description, version, buildFilesDir)
	if err != nil {
		util.Display(
			os.Stderr, true,
			"rpkgm could not add the package %s to the repo. Error: %s",
			name, err,
		)
		os.Exit(1)
	}
}

// Decide decides what to do based on the given strings.
func Decide(repoDB, name, description, version, buildFilesDir, importFile string) {
	// Connect to the database
	dbAdapter, err := database.NewAdapter("sqlite3", repoDB)
	if err != nil {
		util.Display(os.Stderr, true, "rpkgm could not connect to the database. Error: %s", err)
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

	// import a file to the repo
	if importFile != "" {
		ImportPkgs(importFile, dbAdapter)
	}

	// add a package to the repo
	if name != "" && version != "" {
		addPkg(name, description, version, buildFilesDir, dbAdapter)
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
