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
package handlers_test

import (
	"configNexus/internal/handlers"
	"configNexus/internal/utils"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	utils.GlobalDomainPatternsMutex.Lock()
	utils.GlobalDomainPatterns = []utils.RegexPattern{
		{
			Name:  "TestPattern1",
			Regex: "(?P<Function>fn)-(?P<Datacenter>dc)",
		},
	}
	utils.GlobalDomainPatternsMutex.Unlock()
}

func teardown() {
	utils.GlobalDomainPatternsMutex.Lock()
	utils.GlobalDomainPatterns = nil
	utils.GlobalDomainPatternsMutex.Unlock()
}

func TestDetailsHandler(t *testing.T) {
	setup()
	defer teardown()

	h := http.HandlerFunc(handlers.DetailsHandler())

	t.Run("Missing Hostname", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/details/", nil)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		body, _ := io.ReadAll(rr.Body)
		assert.Equal(t, "Missing hostname", strings.TrimSpace(string(body)))
	})

	t.Run("Invalid Regex", func(t *testing.T) {
		// Add an invalid regex pattern
		utils.GlobalDomainPatternsMutex.Lock()
		utils.GlobalDomainPatterns[0].Regex = "(invalid"
		utils.GlobalDomainPatternsMutex.Unlock()

		req := httptest.NewRequest("GET", "/details/test", nil)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		body, _ := io.ReadAll(rr.Body)
		assert.Equal(t, "Invalid Regex Pattern", strings.TrimSpace(string(body)))

		// Revert to valid regex pattern
		utils.GlobalDomainPatternsMutex.Lock()
		utils.GlobalDomainPatterns[0].Regex = "(?P<Function>fn)-(?P<Datacenter>dc)"
		utils.GlobalDomainPatternsMutex.Unlock()
	})

	t.Run("No Matching Pattern", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/details/test", nil)
		rr := httptest.NewRecorder()

		h.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		body, _ := io.ReadAll(rr.Body)
		assert.Equal(t, "No matching pattern found", strings.TrimSpace(string(body)))
	})

	// Add more tests for matching patterns, template processing, etc.
}
