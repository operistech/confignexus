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
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
	"os"
	"time"

	git "github.com/go-git/go-git/v5"
)

var GlobalRepoPath string

// ManageRepo TODO: Add functions for https oauth and for deployment keys
func ManageRepo(repoURL, branch string) error {

	// Create a unique directory within the system's temp folder
	tempDir, err := os.MkdirTemp("", "confignexus")
	if err != nil {
		return err
	}
	log.Debug().Msg(tempDir)

	// Clone the repo
	repo, err := git.PlainClone(tempDir, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		return err
	}

	// Set the global repo path
	GlobalRepoPath = tempDir

	if err := LoadDomainMatchingPatterns(GlobalRepoPath + "/domains_regex.yaml"); err != nil {
		log.Fatal().Err(err).Msg("Failed to load domain matching patterns")
		return err
	}

	// Start a goroutine to pull updates every 20 minutes
	go func() {
		ticker := time.NewTicker(20 * time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C // Wait for the next tick
			wt, err := repo.Worktree()
			if err != nil {
				log.Error().Err(err).Msg("Getting Worktree failed")
			}
			err = wt.Pull(&git.PullOptions{
				RemoteName: "origin",
			})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				log.Error().Err(err).Msg("Failed to pull repository")
			} else {
				// Update domain matches if repo was updated.
				if err := LoadDomainMatchingPatterns(GlobalRepoPath + "/domains_regex.yaml"); err != nil {
					log.Fatal().Err(err).Msg("Failed to load domain matching patterns")
					return
				}
			}

		}
	}()

	return nil
}
