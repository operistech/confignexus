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
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type RegexPattern struct {
	Name  string `yaml:"name"`
	Regex string `yaml:"regex"`
}

type DomainMatching struct {
	RegexPatterns []RegexPattern `yaml:"regex_patterns"`
}

var (
	GlobalDomainPatternsMutex sync.Mutex
	GlobalDomainPatterns      []RegexPattern
)

// LoadDomainMatchingPatterns loads the regex patterns from a YAML file
func LoadDomainMatchingPatterns(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var dm DomainMatching
	if err := yaml.Unmarshal(data, &dm); err != nil {
		return err
	}

	GlobalDomainPatternsMutex.Lock()
	GlobalDomainPatterns = dm.RegexPatterns
	GlobalDomainPatternsMutex.Unlock()

	return nil
}

// GetDomainPatterns safely returns a copy of the global domain patterns
func GetDomainPatterns() []RegexPattern {
	GlobalDomainPatternsMutex.Lock()
	defer GlobalDomainPatternsMutex.Unlock()

	// Create a copy to avoid external modification
	copyPatterns := make([]RegexPattern, len(GlobalDomainPatterns))
	copy(copyPatterns, GlobalDomainPatterns)

	return copyPatterns
}
