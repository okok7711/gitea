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
	tplCodeFrequency base.TplName = "repo/activity"
)

// CodeFrequency renders the page to show repository code frequency
func CodeFrequency(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("repo.activity.navbar.code_frequency")

	ctx.Data["PageIsActivity"] = true
	ctx.Data["PageIsCodeFrequency"] = true
	ctx.PageData["repoLink"] = ctx.Repo.RepoLink

	ctx.HTML(http.StatusOK, tplCodeFrequency)
}

// CodeFrequencyData returns JSON of code frequency data
func CodeFrequencyData(ctx *context.Context) {
	if contributorStats, err := contributors_service.GetContributorStats(ctx, ctx.Cache, ctx.Repo.Repository, ctx.Repo.CommitID); err != nil {
		if errors.Is(err, contributors_service.ErrAwaitGeneration) {
			ctx.Status(http.StatusAccepted)
			return
		}
		ctx.ServerError("GetCodeFrequencyData", err)
	} else {
		ctx.JSON(http.StatusOK, contributorStats["total"].Weeks)
	}
}
