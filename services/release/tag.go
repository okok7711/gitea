// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package release

import (
	"context"
	"errors"
	"fmt"

	"github.com/okok7711/gitea/models/db"
	repo_model "github.com/okok7711/gitea/models/repo"
	"github.com/okok7711/gitea/modules/graceful"
	"github.com/okok7711/gitea/modules/log"
	"github.com/okok7711/gitea/modules/queue"
	repo_module "github.com/okok7711/gitea/modules/repository"

	"xorm.io/builder"
)

type TagSyncOptions struct {
	RepoID int64
}

// tagSyncQueue represents a queue to handle tag sync jobs.
var tagSyncQueue *queue.WorkerPoolQueue[*TagSyncOptions]

func handlerTagSync(items ...*TagSyncOptions) []*TagSyncOptions {
	for _, opts := range items {
		err := repo_module.SyncRepoTags(graceful.GetManager().ShutdownContext(), opts.RepoID)
		if err != nil {
			log.Error("syncRepoTags [%d] failed: %v", opts.RepoID, err)
		}
	}
	return nil
}

func addRepoToTagSyncQueue(repoID int64) error {
	return tagSyncQueue.Push(&TagSyncOptions{
		RepoID: repoID,
	})
}

func initTagSyncQueue(ctx context.Context) error {
	tagSyncQueue = queue.CreateUniqueQueue(ctx, "tag_sync", handlerTagSync)
	if tagSyncQueue == nil {
		return errors.New("unable to create tag_sync queue")
	}
	go graceful.GetManager().RunWithCancel(tagSyncQueue)

	return nil
}

func AddAllRepoTagsToSyncQueue(ctx context.Context) error {
	if err := db.Iterate(ctx, builder.Eq{"is_empty": false}, func(ctx context.Context, repo *repo_model.Repository) error {
		return addRepoToTagSyncQueue(repo.ID)
	}); err != nil {
		return fmt.Errorf("run sync all tags failed: %v", err)
	}
	return nil
}
