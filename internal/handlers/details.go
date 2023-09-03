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
	"configNexus/internal/utils"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// DetailsHandler function
func DetailsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch hostname from the URL path
		hostname := strings.TrimPrefix(r.URL.Path, "/details/")
		if hostname == "" {
			http.Error(w, "Missing hostname", http.StatusBadRequest)
			return
		}

		// Fetch domain patterns using GetDomainPatterns
		domainPatterns := utils.GetDomainPatterns()

		// Loop through domain patterns and match
		for _, pattern := range domainPatterns {
			re, err := regexp.Compile(pattern.Regex)
			if err != nil {
				http.Error(w, "Invalid Regex Pattern", http.StatusInternalServerError)
				return
			}

			match := re.FindStringSubmatch(hostname)
			if match != nil {
				// Create a map to hold the named matched content
				names := re.SubexpNames()
				mapped := make(map[string]string)
				for i, name := range names {
					if i != 0 && name != "" {
						mapped[name] = match[i]
					}
				}

				// Process the main template
				mainTemplate, err := utils.ProcessTemplate(utils.GlobalRepoPath+"/all.yaml", mapped)
				if err != nil {
					http.Error(w, "Failed to process main template", http.StatusInternalServerError)
					log.Error().Err(err).Msg("Failed to process main template")
					return
				}

				// Process the function-specific template
				functionTemplatePath := utils.GlobalRepoPath + "/functions/" + mapped["Function"] + ".yaml"
				if _, err := os.Stat(functionTemplatePath); err == nil {
					functionTemplate, err := utils.ProcessTemplate(functionTemplatePath, mapped)
					if err != nil {
						http.Error(w, "Failed to process function-specific template", http.StatusInternalServerError)
						log.Error().Err(err).Msg("Failed to process function-specific template")
						return
					}

					// Merge mainTemplate and functionTemplate
					for key, value := range functionTemplate {
						mainTemplate[key] = value
					}
				}

				// Process the datacenter-specific template
				datacenterTemplatePath := utils.GlobalRepoPath + "/datacenters/" + mapped["Datacenter"] + ".yaml"
				if _, err := os.Stat(datacenterTemplatePath); err == nil {
					datacenterTemplate, err := utils.ProcessTemplate(datacenterTemplatePath, mapped)
					if err != nil {
						http.Error(w, "Failed to process datacenter-specific template", http.StatusInternalServerError)
						log.Error().Err(err).Msg("Failed to process datacenter-specific template")
						return
					}

					// Merge mainTemplate and datacenterTemplate
					for key, value := range datacenterTemplate {
						mainTemplate[key] = value
					}
				}

				// Process device specific templates
				deviceTemplatePath := utils.GlobalRepoPath + "/devices/" + hostname + ".yaml"
				if _, err := os.Stat(deviceTemplatePath); err == nil {
					deviceTemplate, err := utils.ProcessTemplate(deviceTemplatePath, mapped)
					if err != nil {
						http.Error(w, "Failed to process Device-specific template", http.StatusInternalServerError)
						log.Error().Err(err).Msg("Failed to process Device-specific template")
						return
					}
					// Merge mainTemplate and datacenterTemplate
					for key, value := range deviceTemplate {
						mainTemplate[key] = value
					}

				}
				// Convert merged map to JSON and send it as a response
				jsonData, err := json.Marshal(mainTemplate)
				if err != nil {
					http.Error(w, "Failed to convert to JSON", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonData)
				return
			}
		}

		http.Error(w, "No matching pattern found", http.StatusNotFound)
	}
}
