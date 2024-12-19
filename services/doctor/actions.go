// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package doctor

import (
	"context"
	"fmt"

	"github.com/okok7711/gitea/models/db"
	repo_model "github.com/okok7711/gitea/models/repo"
	unit_model "github.com/okok7711/gitea/models/unit"
	"github.com/okok7711/gitea/modules/log"
	"github.com/okok7711/gitea/modules/optional"
	repo_service "github.com/okok7711/gitea/services/repository"
)

func disableMirrorActionsUnit(ctx context.Context, logger log.Logger, autofix bool) error {
	var reposToFix []*repo_model.Repository

	for page := 1; ; page++ {
		repos, _, err := repo_model.SearchRepository(ctx, &repo_model.SearchRepoOptions{
			ListOptions: db.ListOptions{
				PageSize: repo_model.RepositoryListDefaultPageSize,
				Page:     page,
			},
			Mirror: optional.Some(true),
		})
		if err != nil {
			return fmt.Errorf("SearchRepository: %w", err)
		}
		if len(repos) == 0 {
			break
		}

		for _, repo := range repos {
			if repo.UnitEnabled(ctx, unit_model.TypeActions) {
				reposToFix = append(reposToFix, repo)
			}
		}
	}

	if len(reposToFix) == 0 {
		logger.Info("Found no mirror with actions unit enabled")
	} else {
		logger.Warn("Found %d mirrors with actions unit enabled", len(reposToFix))
	}
	if !autofix || len(reposToFix) == 0 {
		return nil
	}

	for _, repo := range reposToFix {
		if err := repo_service.UpdateRepositoryUnits(ctx, repo, nil, []unit_model.Type{unit_model.TypeActions}); err != nil {
			return err
		}
	}
	logger.Info("Fixed %d mirrors with actions unit enabled", len(reposToFix))

	return nil
}

func init() {
	Register(&Check{
		Title:     "Disable the actions unit for all mirrors",
		Name:      "disable-mirror-actions-unit",
		IsDefault: false,
		Run:       disableMirrorActionsUnit,
		Priority:  9,
	})
}
