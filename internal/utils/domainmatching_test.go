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
	"io/ioutil"
	"os"
	"testing"
)

func createTempYAMLFile(content string) (string, error) {
	tempFile, err := ioutil.TempFile("", "*.yaml")
	if err != nil {
		return "", err
	}

	_, err = tempFile.WriteString(content)
	if err != nil {
		return "", err
	}

	err = tempFile.Close()
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func TestLoadDomainMatchingPatterns(t *testing.T) {
	yamlContent := `
regex_patterns:
  - name: "test"
    regex: ".*"
  - name: "example"
    regex: "^example$"
`

	tempFile, err := createTempYAMLFile(yamlContent)
	if err != nil {
		t.Fatalf("Failed to create temporary YAML file: %v", err)
	}
	defer os.Remove(tempFile)

	t.Run("successfully load patterns from YAML file", func(t *testing.T) {
		err := LoadDomainMatchingPatterns(tempFile)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		patterns := GetDomainPatterns()
		if len(patterns) != 2 {
			t.Fatalf("Expected 2 patterns, got %d", len(patterns))
		}

		expectedNames := []string{"test", "example"}
		expectedRegexes := []string{".*", "^example$"}

		for i, pattern := range patterns {
			if pattern.Name != expectedNames[i] || pattern.Regex != expectedRegexes[i] {
				t.Errorf("Unexpected pattern: got %+v, expected {Name: %s, Regex: %s}", pattern, expectedNames[i], expectedRegexes[i])
			}
		}
	})
}

func TestGetDomainPatterns(t *testing.T) {
	t.Run("fetch patterns without loading should be empty", func(t *testing.T) {
		GlobalDomainPatterns = []RegexPattern{}
		patterns := GetDomainPatterns()

		if len(patterns) != 0 {
			t.Errorf("Expected an empty slice, got: %v", patterns)
		}
	})
}
