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

package main

import (
	"configNexus/internal/handlers"
	"configNexus/internal/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// Clean temporary files
func cleanup() {
	log.Info().Msg("Cleaning Up before Exiting")
	err := os.RemoveAll(utils.GlobalRepoPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to remove temporary folder")
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	// Catch ^C and try to cleanup tmp files on exit
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
	}()

	settings, err := utils.LoadSettings()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load settings")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if settings.DebugLog != "false" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	httpEnabled := settings.HTTPEnabled != "false"
	httpsRedirect := settings.HTTPRedirect != "false"

	defer cleanup()
	err = utils.ManageRepo(settings.RepoAddress, settings.RepoBranch)
	if err != nil {
		log.Fatal().Err(err).Msg("There was a problem with the repo")
	}

	for _, pattern := range utils.GetDomainPatterns() {
		log.Debug().Str(pattern.Name, pattern.Regex).Msg("Regexes")
	}

	// Set up the same handlers for HTTPS
	httpsMux := handlers.SetupHandlers()

	if httpEnabled {
		go func() {
			var httpMux *http.ServeMux
			if httpsRedirect {
				httpMux = http.NewServeMux()
				httpMux.HandleFunc("/", utils.Redirect(settings.ListenAddress, settings.HTTPSPort)) // Redirect to HTTPS
			} else {
				httpMux = httpsMux // Use the same handlers as HTTPS
			}

			log.Info().
				Str("Port", settings.HTTPPort).
				Msg("HTTP server running")
			if err := http.ListenAndServe(settings.HTTPAddr, httpMux); err != nil {
				log.Fatal().Err(err).Msg("HTTP startup failed")
			}
		}()
	}

	// Generate self-signed certificates if they don't exist
	if err := utils.GenerateSelfSignedCert(settings.CertPath, settings.KeyPath); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate self-signed certificates")
	}

	log.Info().
		Str("Port", settings.HTTPSPort).
		Msg("HTTPS server running")
	if err := http.ListenAndServeTLS(settings.HTTPSAddr, settings.CertPath, settings.KeyPath, httpsMux); err != nil {
		log.Fatal().Err(err).Msg("HTTPS startup failed")
	}
}
