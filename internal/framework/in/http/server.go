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

package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/redds-be/redd-go-template/internal/ports"
)

// Adapter is a rpc adapter the is compatible with the api.
type Adapter struct {
	api ports.APIPort
}

// NewAdapter creates a grpc entry adapter that is compatible with gRPCEntryPort.
func NewAdapter(api ports.APIPort) *Adapter {
	return &Adapter{api: api}
}

// Run is the HTTP entrypoint for the application.
func (httpa Adapter) Run() error {
	http.HandleFunc("/", httpa.hello)

	// Listen and serve
	err := http.ListenAndServe(fmt.Sprintf(":%s", "8080"), nil) //nolint:gosec

	return err
}

// hello tells hello world to a http client.
func (httpa Adapter) hello(writer http.ResponseWriter, _ *http.Request) {
	message, err := httpa.api.GetHelloWorld()
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintf(writer, "<h1>%s</h1>", message)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}
}
