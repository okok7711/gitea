// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package issue

import (
	"context"
	"fmt"

	"github.com/okok7711/gitea/models/db"
	issues_model "github.com/okok7711/gitea/models/issues"
	access_model "github.com/okok7711/gitea/models/perm/access"
	repo_model "github.com/okok7711/gitea/models/repo"
	user_model "github.com/okok7711/gitea/models/user"
	"github.com/okok7711/gitea/modules/timeutil"
	notify_service "github.com/okok7711/gitea/services/notify"
)

// CreateRefComment creates a commit reference comment to issue.
func CreateRefComment(ctx context.Context, doer *user_model.User, repo *repo_model.Repository, issue *issues_model.Issue, content, commitSHA string) error {
	if len(commitSHA) == 0 {
		return fmt.Errorf("cannot create reference with empty commit SHA")
	}

	if user_model.IsUserBlockedBy(ctx, doer, issue.PosterID, repo.OwnerID) {
		if isAdmin, _ := access_model.IsUserRepoAdmin(ctx, repo, doer); !isAdmin {
			return user_model.ErrBlockedUser
		}
	}

	// Check if same reference from same commit has already existed.
	has, err := db.GetEngine(ctx).Get(&issues_model.Comment{
		Type:      issues_model.CommentTypeCommitRef,
		IssueID:   issue.ID,
		CommitSHA: commitSHA,
	})
	if err != nil {
		return fmt.Errorf("check reference comment: %w", err)
	} else if has {
		return nil
	}

	_, err = issues_model.CreateComment(ctx, &issues_model.CreateCommentOptions{
		Type:      issues_model.CommentTypeCommitRef,
		Doer:      doer,
		Repo:      repo,
		Issue:     issue,
		CommitSHA: commitSHA,
		Content:   content,
	})
	return err
}

// CreateIssueComment creates a plain issue comment.
func CreateIssueComment(ctx context.Context, doer *user_model.User, repo *repo_model.Repository, issue *issues_model.Issue, content string, attachments []string) (*issues_model.Comment, error) {
	if user_model.IsUserBlockedBy(ctx, doer, issue.PosterID, repo.OwnerID) {
		if isAdmin, _ := access_model.IsUserRepoAdmin(ctx, repo, doer); !isAdmin {
			return nil, user_model.ErrBlockedUser
		}
	}

	comment, err := issues_model.CreateComment(ctx, &issues_model.CreateCommentOptions{
		Type:        issues_model.CommentTypeComment,
		Doer:        doer,
		Repo:        repo,
		Issue:       issue,
		Content:     content,
		Attachments: attachments,
	})
	if err != nil {
		return nil, err
	}

	mentions, err := issues_model.FindAndUpdateIssueMentions(ctx, issue, doer, comment.Content)
	if err != nil {
		return nil, err
	}

	notify_service.CreateIssueComment(ctx, doer, repo, issue, comment, mentions)

	return comment, nil
}

// UpdateComment updates information of comment.
func UpdateComment(ctx context.Context, c *issues_model.Comment, contentVersion int, doer *user_model.User, oldContent string) error {
	if err := c.LoadIssue(ctx); err != nil {
		return err
	}
	if err := c.Issue.LoadRepo(ctx); err != nil {
		return err
	}

	if user_model.IsUserBlockedBy(ctx, doer, c.Issue.PosterID, c.Issue.Repo.OwnerID) {
		if isAdmin, _ := access_model.IsUserRepoAdmin(ctx, c.Issue.Repo, doer); !isAdmin {
			return user_model.ErrBlockedUser
		}
	}

	needsContentHistory := c.Content != oldContent && c.Type.HasContentSupport()
	if needsContentHistory {
		hasContentHistory, err := issues_model.HasIssueContentHistory(ctx, c.IssueID, c.ID)
		if err != nil {
			return err
		}
		if !hasContentHistory {
			if err = issues_model.SaveIssueContentHistory(ctx, c.PosterID, c.IssueID, c.ID,
				c.CreatedUnix, oldContent, true); err != nil {
				return err
			}
		}
	}

	if err := issues_model.UpdateComment(ctx, c, contentVersion, doer); err != nil {
		return err
	}

	if needsContentHistory {
		err := issues_model.SaveIssueContentHistory(ctx, doer.ID, c.IssueID, c.ID, timeutil.TimeStampNow(), c.Content, false)
		if err != nil {
			return err
		}
	}

	notify_service.UpdateComment(ctx, doer, c, oldContent)

	return nil
}

// DeleteComment deletes the comment
func DeleteComment(ctx context.Context, doer *user_model.User, comment *issues_model.Comment) error {
	err := db.WithTx(ctx, func(ctx context.Context) error {
		return issues_model.DeleteComment(ctx, comment)
	})
	if err != nil {
		return err
	}

	notify_service.DeleteComment(ctx, doer, comment)

	return nil
}
