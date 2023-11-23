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
	"net/http/httptest"
	"testing"
)

func TestRedirect(t *testing.T) {
	tests := []struct {
		name           string
		listenAddress  string
		httpsPort      string
		requestPath    string
		expectedStatus int
		expectedURL    string
	}{
		{
			name:           "standard port",
			listenAddress:  "example.com",
			httpsPort:      "443",
			requestPath:    "/somepath",
			expectedStatus: http.StatusMovedPermanently,
			expectedURL:    "https://example.com/somepath",
		},
		{
			name:           "non-standard port",
			listenAddress:  "example.com",
			httpsPort:      "8443",
			requestPath:    "/anotherpath",
			expectedStatus: http.StatusMovedPermanently,
			expectedURL:    "https://example.com:8443/anotherpath",
		},
		{
			name:           "empty path",
			listenAddress:  "example.com",
			httpsPort:      "443",
			requestPath:    "",
			expectedStatus: http.StatusMovedPermanently,
			expectedURL:    "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://"+tt.listenAddress+tt.requestPath, nil)
			w := httptest.NewRecorder()

			handler := Redirect(tt.listenAddress, tt.httpsPort)
			handler(w, req)

			resp := w.Result()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, resp.StatusCode)
			}

			location, ok := resp.Header["Location"]
			if !ok {
				t.Errorf("Location header not set")
			} else if location[0] != tt.expectedURL {
				t.Errorf("Expected URL %s; got %s", tt.expectedURL, location[0])
			}
		})
	}
}
