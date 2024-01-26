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

package api

import (
	"log"
	"os"
	"testing"

	"github.com/redds-be/redd-go-template/internal/application/core/helloworld"
	"github.com/redds-be/redd-go-template/internal/framework/out/database"
	"github.com/stretchr/testify/suite"
)

func (s *APISuite) TestGetHelloWorld() {
	message, err := s.testApp.GetHelloWorld()
	s.Require().NoError(err)
	s.Equal("Hello, World!", message)
}

type APISuite struct {
	testApp Application
	suite.Suite
}

func TestAPISuite(t *testing.T) {
	dbAdapter, err := database.NewAdapter("sqlite3", "test_api.db")
	if err != nil {
		log.Fatal(err)
	}

	err = dbAdapter.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	defer func(dbAdapter *database.Adapter) {
		err := dbAdapter.CloseDBConnection()
		if err != nil {
			t.Fatal("unable to close the database")
		}
	}(dbAdapter)

	defer func() {
		err := os.Remove("test_api.db")
		if err != nil {
			t.Fatal("unable to remove the test database")
		}
	}()

	core := helloworld.New()

	t.Parallel()
	suite.Run(t, &APISuite{
		testApp: Application{dbase: dbAdapter, hello: core},
		Suite:   suite.Suite{},
	})
}
