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
    buildFilesDir VARCHAR(4096) NOT NULL
    );`
	_, err := dbAdapter.dbase.Exec(queryString)

	return err
}

// AddToMainRepo adds a package to the package table in the repo.
func (dbAdapter Adapter) AddToMainRepo(name, description, repoVersion, buildFilesDir string) error {
	const queryString = `INSERT INTO packages VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := dbAdapter.dbase.Exec(
		queryString,
		name,
		description,
		repoVersion,
		"",
		false,
		buildFilesDir,
	)

	return err
}

// GetPkgInfo returns the basic information about a given package.
func (dbAdapter Adapter) GetPkgInfo(name string) (PkgInfo, error) {
	const queryString = `SELECT name, description, repoVersion, installedVersion, installed FROM packages WHERE name = $1;`
	var info PkgInfo

	err := dbAdapter.dbase.QueryRow(queryString, name).
		Scan(&info.Name, &info.Description, &info.RepoVersion, &info.InstalledVersion, &info.Installed)
	if err != nil {
		return PkgInfo{}, err
	}

	return info, err
}

// GetAllPkgInfo returns the basic information about all packages in the repo.
func (dbAdapter Adapter) GetAllPkgInfo() ([]PkgInfo, error) {
	const queryString = `SELECT name, description, repoVersion, installedVersion, installed FROM packages;`
	var info []PkgInfo

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
		var i PkgInfo
		err = rows.Scan(&i.Name, &i.Description, &i.RepoVersion, &i.InstalledVersion, &i.Installed)
		info = append(info, i)
	}

	return info, err
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
