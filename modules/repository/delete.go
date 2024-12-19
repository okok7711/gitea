// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repository

import (
	"context"

	"github.com/okok7711/gitea/models/organization"
	repo_model "github.com/okok7711/gitea/models/repo"
	user_model "github.com/okok7711/gitea/models/user"
)

// CanUserDelete returns true if user could delete the repository
func CanUserDelete(ctx context.Context, repo *repo_model.Repository, user *user_model.User) (bool, error) {
	if user.IsAdmin || user.ID == repo.OwnerID {
		return true, nil
	}

	if err := repo.LoadOwner(ctx); err != nil {
		return false, err
	}

	if repo.Owner.IsOrganization() {
		isAdmin, err := organization.OrgFromUser(repo.Owner).IsOrgAdmin(ctx, user.ID)
		if err != nil {
			return false, err
		}
		return isAdmin, nil
	}

	return false, nil
}
