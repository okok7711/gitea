// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package feed

import (
	"github.com/okok7711/gitea/services/context"
)

// RenderBranchFeed render format for branch or file
func RenderBranchFeed(ctx *context.Context) {
	_, _, showFeedType := GetFeedType(ctx.PathParam(":reponame"), ctx.Req)
	if ctx.Repo.TreePath == "" {
		ShowBranchFeed(ctx, ctx.Repo.Repository, showFeedType)
	} else {
		ShowFileFeed(ctx, ctx.Repo.Repository, showFeedType)
	}
}
