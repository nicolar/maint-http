// maint-http: fast-superlight http server for maintenance pages
// Copyright (c) Nicola Ruggero 2020 <nicola@nxnt.org>
//
// This is a fast-superlight http server for maintenance pages,
// specifically designed to be a "backup" backend for haproxy.
// For any request a maintenance page is served.
//
// Usage: ./maint-http [-d <directory_to_serve>] [-p <tcp_listen_port>] [-r <html_resources_context>]
//
// License: GPL v3+
// ====================================================================
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
// ====================================================================

// Version 1.0 - 12/04/2020 - Initial version
// Version 1.1 - 20/04/2020 - Sanitization, resource context customizable

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Set globals
var httpPort *string
var dir *string
var ctxRes *string
var showSwVer *bool

const swVer = "1.1"

// Special handler to always serve a maintenance page and its resources
func handlerMaint(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Printf("Before: %s", r.URL.String())

		// Detect html_resource_context requests
		re := regexp.MustCompile(*ctxRes)

		// Strip unwanted prefix for html_resource_context requests
		reCapture := regexp.MustCompile(*ctxRes + ".*")

		if re.MatchString(r.URL.Path) {
			r.URL.Path = reCapture.FindString(r.URL.Path)
		} else {
			r.URL.Path = "/"
		}
		h.ServeHTTP(w, r) // call original
		//log.Printf("After: %s", r.URL.String())
	})
}

func main() {

	// Options
	httpPort = flag.String("p", "3000", "TCP listen port")
	dir = flag.String("d", "./", "Directory with maintenance page to serve")
	ctxRes = flag.String("r", "res", "URL context for HTML resources")
	showSwVer = flag.Bool("v", false, "Print software version and exit")
	flag.Parse()

	// Show Software version
	if *showSwVer {
		fmt.Printf("maint-http: fast-superlight http server for maintenance pages\n")
		fmt.Printf("Version: %s\n", swVer)
		os.Exit(1)
	}

	// Sanitize html_resource_context
	*ctxRes = "/" + strings.TrimPrefix(strings.TrimSuffix(*ctxRes, "/"), "/") + "/"

	// Sanitize directory_to_serve
	*dir = strings.TrimSuffix(*dir, "/") + "/"
	// Standard FileServer handler for maintenance page
	fs := http.FileServer(http.Dir(*dir))

	http.Handle("/", handlerMaint(fs))

	log.Printf("Starting maint-http: fast-superlight http server for maintenance pages\n")
	log.Printf("Version: %s\n", swVer)
	log.Printf("Serving %s (resources: %s) on http://localhost:%s\n", *dir, *ctxRes, *httpPort)
	err := http.ListenAndServe("localhost:"+*httpPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
