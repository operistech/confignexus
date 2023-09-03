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
	"bytes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"os"
	"text/template"
)

func ProcessTemplate(filePath string, data map[string]string) (map[string]interface{}, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		for key, value := range data {
			log.Debug().Str(key, value).Msg("Template data map")
		}
	}

	tmpl, err := template.New("config").Parse(string(fileContent))
	if err != nil {
		return nil, err
	}

	// Create a buffer to store the output
	var output bytes.Buffer

	// Execute the template
	err = tmpl.Execute(&output, data)
	if err != nil {
		return nil, err
	}

	// Assuming that the template generates YAML content,
	// we need to Unmarshal it into a map for further use.
	var yamlMap map[string]interface{}
	err = yaml.Unmarshal(output.Bytes(), &yamlMap)
	if err != nil {
		return nil, err
	}

	return yamlMap, nil
}
