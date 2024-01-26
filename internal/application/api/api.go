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

import "github.com/redds-be/redd-go-template/internal/ports"

type Application struct {
	dbase ports.DBPort
	hello Helloworld
}

// NewApplication creates a new Application.
func NewApplication(dbase ports.DBPort, hello Helloworld) *Application {
	return &Application{dbase: dbase, hello: hello}
}

// GetHelloWorld gets the hello world message.
func (apia Application) GetHelloWorld() (string, error) {
	message := apia.hello.HelloWorld()

	err := apia.dbase.AddToHistory(message)
	if err != nil {
		return "", err
	}

	return message, nil
}
