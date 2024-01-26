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

package helloworld

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func (s *HelloSuite) TestHello() {
	adapter := New()
	message := adapter.HelloWorld()
	s.Equal("Hello, World!", message)
}

type HelloSuite struct {
	suite.Suite
}

func TestHelloSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HelloSuite))
}
