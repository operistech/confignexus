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
package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupHandlers(t *testing.T) {
	mux := SetupHandlers()

	t.Run("root endpoint returns welcome message", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status OK; got %v", resp.Status)
		}

		expectedBody := "Welcome to ConfigNexus!"
		if string(body) != expectedBody {
			t.Errorf("Expected body to be '%s'; got '%s'", expectedBody, body)
		}
	})

	t.Run("/details/ endpoint is handled", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/details/", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode == http.StatusNotFound {
			t.Errorf("Expected /details/ to be handled; got status %v", resp.Status)
		}
	})
}
