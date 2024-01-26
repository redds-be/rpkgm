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
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func (s *DBSuite) TestDB() {
	// Connect to the database
	frameRight, err := NewAdapter("sqlite3", "test_database.db")
	s.Require().NoError(err)

	// Create a table
	err = frameRight.CreateTable()
	s.Require().NoError(err)

	// Add an entry to the db
	err = frameRight.AddToHistory("Hello, World!")
	s.Require().NoError(err)

	// Close database connection
	err = frameRight.CloseDBConnection()
	s.Require().NoError(err)

	// Remove the database
	err = os.Remove("test_database.db")
	s.Require().NoError(err)
}

type DBSuite struct {
	suite.Suite
}

func TestDBSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DBSuite))
}
