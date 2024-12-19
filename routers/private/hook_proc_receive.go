// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package private

import (
	"net/http"

	repo_model "github.com/okok7711/gitea/models/repo"
	"github.com/okok7711/gitea/modules/git"
	"github.com/okok7711/gitea/modules/log"
	"github.com/okok7711/gitea/modules/private"
	"github.com/okok7711/gitea/modules/web"
	"github.com/okok7711/gitea/services/agit"
	gitea_context "github.com/okok7711/gitea/services/context"
)

// HookProcReceive proc-receive hook - only handles agit Proc-Receive requests at present
func HookProcReceive(ctx *gitea_context.PrivateContext) {
	opts := web.GetForm(ctx).(*private.HookOptions)
	if !git.DefaultFeatures().SupportProcReceive {
		ctx.Status(http.StatusNotFound)
		return
	}

	results, err := agit.ProcReceive(ctx, ctx.Repo.Repository, ctx.Repo.GitRepo, opts)
	if err != nil {
		if repo_model.IsErrUserDoesNotHaveAccessToRepo(err) {
			ctx.Error(http.StatusBadRequest, "UserDoesNotHaveAccessToRepo", err.Error())
		} else {
			log.Error(err.Error())
			ctx.JSON(http.StatusInternalServerError, private.Response{
				Err: err.Error(),
			})
		}

		return
	}

	ctx.JSON(http.StatusOK, private.HookProcReceiveResult{
		Results: results,
	})
}
