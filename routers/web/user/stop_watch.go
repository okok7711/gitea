// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"net/http"

	"github.com/okok7711/gitea/models/db"
	issues_model "github.com/okok7711/gitea/models/issues"
	"github.com/okok7711/gitea/services/context"
	"github.com/okok7711/gitea/services/convert"
)

// GetStopwatches get all stopwatches
func GetStopwatches(ctx *context.Context) {
	sws, err := issues_model.GetUserStopwatches(ctx, ctx.Doer.ID, db.ListOptions{
		Page:     ctx.FormInt("page"),
		PageSize: convert.ToCorrectPageSize(ctx.FormInt("limit")),
	})
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	count, err := issues_model.CountUserStopwatches(ctx, ctx.Doer.ID)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	apiSWs, err := convert.ToStopWatches(ctx, sws)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetTotalCountHeader(count)
	ctx.JSON(http.StatusOK, apiSWs)
}
