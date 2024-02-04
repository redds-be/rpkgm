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

package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // Driver for sqlite
)

// Package defines a package in the database.
type Package struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Version       string `json:"version"`
	BuildFilesDir string `json:"buildFilesDir"`
	ArchiveURL    string `json:"archiveUrl"`
	Sha512        string `json:"sha512"`
	Dependencies  string `json:"dependencies"`
}

// Packages defines a slice of package.
type Packages struct {
	Packages []Package `json:"packages"`
}

// PkgInfo defines the basic information about a give package.
type PkgInfo struct {
	Name             string
	Description      string
	RepoVersion      string
	InstalledVersion string
	Installed        bool
	BuildFilesDir    string
	ArchiveURL       string
	Sha512           string
	Dependencies     string
}

// Adapter implements the DBPort interface.
type Adapter struct {
	dbase *sql.DB
}

// NewAdapter creates a new Adapter.
func NewAdapter(driverName, dataSourceName string) (*Adapter, error) {
	// Connect to the database
	dbase, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// Test db connection
	err = dbase.Ping()
	if err != nil {
		return nil, err
	}

	return &Adapter{dbase: dbase}, nil
}

// CloseDBConnection closes the db connection.
func (dbAdapter Adapter) CloseDBConnection() error {
	// Close the database
	err := dbAdapter.dbase.Close()

	return err
}

// CreatePkgTable create the packages table.
func (dbAdapter Adapter) CreatePkgTable() error {
	const queryString = `CREATE TABLE IF NOT EXISTS packages (
    name VARCHAR(512) PRIMARY KEY,
    description VARCHAR(8000) NOT NULL,
    repoVersion VARCHAR(16) NOT NULL,
    installedVersion VARCHAR(16),
    installed BOOLEAN NOT NULL,
    buildFilesDir VARCHAR(4096) NOT NULL,
    archiveURL VARCHAR(8000) NOT NULL,
    sha512 VARCHAR(128) NOT NULL,
    dependencies VARCHAR(8000) NOT NULL
    );`
	_, err := dbAdapter.dbase.Exec(queryString)

	return err
}

// AddToRepo adds a package to the package table in the repo.
func (dbAdapter Adapter) AddToRepo(
	name, description, repoVersion, buildFilesDir, archiveURL, hash, dependencies string,
) error {
	const queryString = `INSERT INTO packages VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`
	_, err := dbAdapter.dbase.Exec(
		queryString,
		name,
		description,
		repoVersion,
		"",
		false,
		buildFilesDir,
		archiveURL,
		hash,
		dependencies,
	)

	return err
}

// SyncRepo syncs packages in the database (using the name as the key).
func (dbAdapter Adapter) SyncRepo(
	name, description, repoVersion, buildFilesDir, archiveURL, hash, dependencies string,
) error {
	const queryString = `UPDATE packages SET
        description = $1,
        repoVersion = $2,
        buildFilesDir = $3,
        archiveURL = $4,
        sha512 = $5,
        dependencies = $6 
        WHERE name = $7;`

	_, err := dbAdapter.dbase.Exec(
		queryString,
		description,
		repoVersion,
		buildFilesDir,
		archiveURL,
		hash,
		dependencies,
		name,
	)

	return err
}

// GetPkgInfo returns the basic information about a given package.
func (dbAdapter Adapter) GetPkgInfo(name string) (PkgInfo, error) {
	const queryString = `SELECT
        name,
        description,
        repoVersion,
        installedVersion,
        installed,
        buildFilesDir,
        archiveURL,
        sha512,
        dependencies
        FROM packages WHERE name = $1;`

	var info PkgInfo

	err := dbAdapter.dbase.QueryRow(queryString, name).Scan(
		&info.Name,
		&info.Description,
		&info.RepoVersion,
		&info.InstalledVersion,
		&info.Installed,
		&info.BuildFilesDir,
		&info.ArchiveURL,
		&info.Sha512,
		&info.Dependencies,
	)
	if err != nil {
		return PkgInfo{}, err
	}

	return info, err
}

// GetAllPkgInfo returns the basic information about all packages in the repo.
func (dbAdapter Adapter) GetAllPkgInfo() ([]PkgInfo, error) {
	const queryString = `SELECT
        name,
        description,
        repoVersion,
        installedVersion,
        installed,
        buildFilesDir,
        archiveURL,
        sha512,
        dependencies
        FROM packages;`

	var infos []PkgInfo

	// Get the row results of the query
	rows, err := dbAdapter.dbase.Query(queryString) //nolint:sqlclosecheck
	if err != nil {
		return nil, err
	}

	// Defer the closing of the rows
	defer func(rows *sql.Rows) {
		err = rows.Close()
	}(rows)

	if rows.Err() != nil {
		return nil, err
	}

	// For each row, append to info
	for rows.Next() {
		var info PkgInfo
		err = rows.Scan(
			&info.Name,
			&info.Description,
			&info.RepoVersion,
			&info.InstalledVersion,
			&info.Installed,
			&info.BuildFilesDir,
			&info.ArchiveURL,
			&info.Sha512,
			&info.Dependencies,
		)
		infos = append(infos, info)
	}

	return infos, err
}

// RenamePackage renames a given package.
func (dbAdapter Adapter) RenamePackage(oldName, newName string) error {
	const queryString = `UPDATE packages SET name = $1 WHERE name = $2;`

	_, err := dbAdapter.dbase.Exec(queryString, newName, oldName)
	if err != nil {
		return err
	}

	return nil
}

// ChangePkgDesc changes a given package's description.
func (dbAdapter Adapter) ChangePkgDesc(name, description string) error {
	const queryString = `UPDATE packages SET description = $1 WHERE name = $2;`
	_, err := dbAdapter.dbase.Exec(queryString, description, name)

	return err
}

// GetPkgBuildFilesDir returns the build files dir for a given package.
func (dbAdapter Adapter) GetPkgBuildFilesDir(name string) (string, error) {
	const queryString = `SELECT buildFilesDir FROM packages WHERE name = $1;`
	var buildFilesDir string

	err := dbAdapter.dbase.QueryRow(queryString, name).Scan(&buildFilesDir)
	if err != nil {
		return "", err
	}

	return buildFilesDir, nil
}

// GetRepoVersion returns the repo's version of a given package from the repo.
func (dbAdapter Adapter) GetRepoVersion(name string) (string, error) {
	const queryString = `SELECT repoVersion FROM packages WHERE name = $1;`
	var version string

	err := dbAdapter.dbase.QueryRow(queryString, name).Scan(&version)
	if err != nil {
		return "", err
	}

	return version, nil
}

// SetInstalledVersion sets the installed version for a package.
func (dbAdapter Adapter) SetInstalledVersion(name, version string) error {
	const queryString = `UPDATE packages SET installedVersion = $1 WHERE name = $2;`
	_, err := dbAdapter.dbase.Exec(queryString, version, name)

	return err
}

// UpdateRepoVersion updates the repo's version of a package.
func (dbAdapter Adapter) UpdateRepoVersion(name, version string) error {
	const queryString = `UPDATE packages SET repoVersion = $1 WHERE name = $2;`
	_, err := dbAdapter.dbase.Exec(queryString, version, name)

	return err
}

// IsPkgInRepo searches the database for a package based on its name.
func (dbAdapter Adapter) IsPkgInRepo(name string) (bool, error) {
	const queryString = `SELECT name FROM packages WHERE name = $1;`
	var nameInDB string

	err := dbAdapter.dbase.QueryRow(queryString, name).Scan(&nameInDB)
	if err != nil {
		return false, err
	}

	return true, nil
}

// IsInstalled checks if a given packages is already installed or not.
func (dbAdapter Adapter) IsInstalled(name string) (bool, error) {
	const queryString = `SELECT installed FROM packages WHERE name = $1;`
	var installed bool

	err := dbAdapter.dbase.QueryRow(queryString, name).Scan(&installed)
	if err != nil {
		return false, err
	}

	return installed, nil
}

// MarkAsInstalled marks a package as installed.
func (dbAdapter Adapter) MarkAsInstalled(name string) error {
	const queryString = `UPDATE packages SET installed = 1 WHERE name = $1;`

	_, err := dbAdapter.dbase.Exec(queryString, name)
	if err != nil {
		return err
	}

	return nil
}

// MarkAsNotInstalled marks a package as installed.
func (dbAdapter Adapter) MarkAsNotInstalled(name string) error {
	const queryString = `UPDATE packages SET installed = 0 WHERE name = $1;`

	_, err := dbAdapter.dbase.Exec(queryString, name)
	if err != nil {
		return err
	}

	return nil
}

// ChangeBuildFilesDir changes the build files directory for a given package.
func (dbAdapter Adapter) ChangeBuildFilesDir(name, buildFilesDir string) error {
	const queryString = `UPDATE packages SET buildFilesDir = $1 WHERE name = $2;`

	_, err := dbAdapter.dbase.Exec(queryString, buildFilesDir, name)
	if err != nil {
		return err
	}

	return nil
}

// RemovePackage removes a given package from the repository.
func (dbAdapter Adapter) RemovePackage(name string) error {
	const queryString = `DELETE FROM packages WHERE name = $1;`

	_, err := dbAdapter.dbase.Exec(queryString, name)
	if err != nil {
		return err
	}

	return nil
}

// ChangeArchiveURL changes the archive url.
func (dbAdapter Adapter) ChangeArchiveURL(name, archiveURL string) error {
	const queryString = `UPDATE packages SET archiveURL = $1 WHERE name = $2;`

	_, err := dbAdapter.dbase.Exec(queryString, archiveURL, name)
	if err != nil {
		return err
	}

	return nil
}

// ChangeHash changes an archive's hash.
func (dbAdapter Adapter) ChangeHash(name, hash string) error {
	const queryString = `UPDATE packages SET sha512 = $1 WHERE name = $2;`

	_, err := dbAdapter.dbase.Exec(queryString, hash, name)
	if err != nil {
		return err
	}

	return nil
}

// ChangeDeps changes a package's dependencies list.
func (dbAdapter Adapter) ChangeDeps(name, dependencies string) error {
	const queryString = `UPDATE packages SET dependencies = $1 WHERE name = $2;`

	_, err := dbAdapter.dbase.Exec(queryString, dependencies, name)
	if err != nil {
		return err
	}

	return nil
}
