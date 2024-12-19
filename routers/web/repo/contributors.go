// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"errors"
	"net/http"

	"github.com/okok7711/gitea/modules/base"
	"github.com/okok7711/gitea/services/context"
	contributors_service "github.com/okok7711/gitea/services/repository"
)

const (
	tplContributors base.TplName = "repo/activity"
)

// Contributors render the page to show repository contributors graph
func Contributors(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.activity.navbar.contributors")
	ctx.Data["PageIsActivity"] = true
	ctx.Data["PageIsContributors"] = true
	ctx.HTML(http.StatusOK, tplContributors)
}

// ContributorsData renders JSON of contributors along with their weekly commit statistics
func ContributorsData(ctx *context.Context) {
	if contributorStats, err := contributors_service.GetContributorStats(ctx, ctx.Cache, ctx.Repo.Repository, ctx.Repo.CommitID); err != nil {
		if errors.Is(err, contributors_service.ErrAwaitGeneration) {
			ctx.Status(http.StatusAccepted)
			return
		}
		ctx.ServerError("GetContributorStats", err)
	} else {
		ctx.JSON(http.StatusOK, contributorStats)
	}
}
