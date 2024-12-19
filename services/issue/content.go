// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package issue

import (
	"context"

	issues_model "github.com/okok7711/gitea/models/issues"
	access_model "github.com/okok7711/gitea/models/perm/access"
	user_model "github.com/okok7711/gitea/models/user"
	notify_service "github.com/okok7711/gitea/services/notify"
)

// ChangeContent changes issue content, as the given user.
func ChangeContent(ctx context.Context, issue *issues_model.Issue, doer *user_model.User, content string, contentVersion int) error {
	if err := issue.LoadRepo(ctx); err != nil {
		return err
	}

	if user_model.IsUserBlockedBy(ctx, doer, issue.PosterID, issue.Repo.OwnerID) {
		if isAdmin, _ := access_model.IsUserRepoAdmin(ctx, issue.Repo, doer); !isAdmin {
			return user_model.ErrBlockedUser
		}
	}

	oldContent := issue.Content

	if err := issues_model.ChangeIssueContent(ctx, issue, doer, content, contentVersion); err != nil {
		return err
	}

	notify_service.IssueChangeContent(ctx, doer, issue, oldContent)

	return nil
}
