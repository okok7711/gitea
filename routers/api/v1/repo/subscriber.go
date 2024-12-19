// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"net/http"

	repo_model "github.com/okok7711/gitea/models/repo"
	api "github.com/okok7711/gitea/modules/structs"
	"github.com/okok7711/gitea/routers/api/v1/utils"
	"github.com/okok7711/gitea/services/context"
	"github.com/okok7711/gitea/services/convert"
)

// ListSubscribers list a repo's subscribers (i.e. watchers)
func ListSubscribers(ctx *context.APIContext) {
	// swagger:operation GET /repos/{owner}/{repo}/subscribers repository repoListSubscribers
	// ---
	// summary: List a repo's watchers
	// produces:
	// - application/json
	// parameters:
	// - name: owner
	//   in: path
	//   description: owner of the repo
	//   type: string
	//   required: true
	// - name: repo
	//   in: path
	//   description: name of the repo
	//   type: string
	//   required: true
	// - name: page
	//   in: query
	//   description: page number of results to return (1-based)
	//   type: integer
	// - name: limit
	//   in: query
	//   description: page size of results
	//   type: integer
	// responses:
	//   "200":
	//     "$ref": "#/responses/UserList"
	//   "404":
	//     "$ref": "#/responses/notFound"

	subscribers, err := repo_model.GetRepoWatchers(ctx, ctx.Repo.Repository.ID, utils.GetListOptions(ctx))
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetRepoWatchers", err)
		return
	}
	users := make([]*api.User, len(subscribers))
	for i, subscriber := range subscribers {
		users[i] = convert.ToUser(ctx, subscriber, ctx.Doer)
	}

	ctx.SetTotalCountHeader(int64(ctx.Repo.Repository.NumWatches))
	ctx.JSON(http.StatusOK, users)
}
