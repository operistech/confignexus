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
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestProcessTemplate(t *testing.T) {
	// Setup temporary files for testing
	tempFile, err := ioutil.TempFile("", "template")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.WriteString("key: {{.key}}\nvalue: {{.value}}")
	tempFile.Close()

	invalidTempFile, err := ioutil.TempFile("", "invalid-template")
	if err != nil {
		t.Fatalf("Failed to create invalid temporary file: %s", err)
	}
	defer os.Remove(invalidTempFile.Name())
	invalidTempFile.WriteString("key: {{.key\nvalue: {{.value}}")
	invalidTempFile.Close()

	tests := []struct {
		name        string
		filePath    string
		data        map[string]string
		expected    map[string]interface{}
		expectedErr bool
		logLevel    zerolog.Level
	}{
		{
			name:        "successful processing",
			filePath:    tempFile.Name(),
			data:        map[string]string{"key": "Key1", "value": "Value1"},
			expected:    map[string]interface{}{"key": "Key1", "value": "Value1"},
			expectedErr: false,
			logLevel:    zerolog.InfoLevel,
		},
		{
			name:        "file not found",
			filePath:    "nonexistentfile.yaml",
			data:        map[string]string{"key": "Key1", "value": "Value1"},
			expected:    nil,
			expectedErr: true,
			logLevel:    zerolog.InfoLevel,
		},
		{
			name:        "invalid template",
			filePath:    invalidTempFile.Name(),
			data:        map[string]string{"key": "Key1", "value": "Value1"},
			expected:    nil,
			expectedErr: true,
			logLevel:    zerolog.InfoLevel,
		},
		{
			name:        "debug level logs",
			filePath:    tempFile.Name(),
			data:        map[string]string{"key": "Key1", "value": "Value1"},
			expected:    map[string]interface{}{"key": "Key1", "value": "Value1"},
			expectedErr: false,
			logLevel:    zerolog.DebugLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global log level
			zerolog.SetGlobalLevel(tt.logLevel)

			got, err := ProcessTemplate(tt.filePath, tt.data)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}

			// Reset to default log level
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		})
	}
}
