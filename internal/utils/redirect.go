/*
   This file is part of configNexus.

   configNexus is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   configNexus is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with configNexus.  If not, see <https://www.gnu.org/licenses/>.

   Copyright (C) 2023 Operistech Inc.
*/

package utils

import (
	"net/http"
)

// Redirect function that redirects HTTP to HTTPS, taking into account the host and port
func Redirect(listenAddress string, httpsPort string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		target := "https://" + listenAddress
		if httpsPort != "443" {
			target += ":" + httpsPort
		}
		target += req.URL.Path
		http.Redirect(w, req, target, http.StatusMovedPermanently)
	}
}
