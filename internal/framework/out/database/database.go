//    redd-go-template, a template for go projects.
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
	"time"

	_ "github.com/mattn/go-sqlite3" // Driver for sqlite
)

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

// CreateTable create the default 'hello_history' table.
func (dbAdapter Adapter) CreateTable() error {
	const queryString = `CREATE TABLE IF NOT EXISTS hello_history (date TIMESTAMP, message VARCHAR(64));`
	_, err := dbAdapter.dbase.Exec(queryString)

	return err
}

// AddToHistory adds the helloworld message to the database history table.
func (dbAdapter Adapter) AddToHistory(message string) error {
	const queryString = `INSERT INTO hello_history VALUES ($1, $2)`
	_, err := dbAdapter.dbase.Exec(queryString, time.Now().UTC(), message)

	return err
}
