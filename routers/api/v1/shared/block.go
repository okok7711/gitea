// Copyright 2024 The Gitea Authors.
// SPDX-License-Identifier: MIT

package shared

import (
	"errors"
	"net/http"

	user_model "github.com/okok7711/gitea/models/user"
	api "github.com/okok7711/gitea/modules/structs"
	"github.com/okok7711/gitea/routers/api/v1/utils"
	"github.com/okok7711/gitea/services/context"
	"github.com/okok7711/gitea/services/convert"
	user_service "github.com/okok7711/gitea/services/user"
)

func ListBlocks(ctx *context.APIContext, blocker *user_model.User) {
	blocks, total, err := user_model.FindBlockings(ctx, &user_model.FindBlockingOptions{
		ListOptions: utils.GetListOptions(ctx),
		BlockerID:   blocker.ID,
	})
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "FindBlockings", err)
		return
	}

	if err := user_model.BlockingList(blocks).LoadAttributes(ctx); err != nil {
		ctx.Error(http.StatusInternalServerError, "LoadAttributes", err)
		return
	}

	users := make([]*api.User, 0, len(blocks))
	for _, b := range blocks {
		users = append(users, convert.ToUser(ctx, b.Blockee, blocker))
	}

	ctx.SetTotalCountHeader(total)
	ctx.JSON(http.StatusOK, &users)
}

func CheckUserBlock(ctx *context.APIContext, blocker *user_model.User) {
	blockee, err := user_model.GetUserByName(ctx, ctx.PathParam("username"))
	if err != nil {
		ctx.NotFound("GetUserByName", err)
		return
	}

	status := http.StatusNotFound
	blocking, err := user_model.GetBlocking(ctx, blocker.ID, blockee.ID)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "GetBlocking", err)
		return
	}
	if blocking != nil {
		status = http.StatusNoContent
	}

	ctx.Status(status)
}

func BlockUser(ctx *context.APIContext, blocker *user_model.User) {
	blockee, err := user_model.GetUserByName(ctx, ctx.PathParam("username"))
	if err != nil {
		ctx.NotFound("GetUserByName", err)
		return
	}

	if err := user_service.BlockUser(ctx, ctx.Doer, blocker, blockee, ctx.FormString("note")); err != nil {
		if errors.Is(err, user_model.ErrCanNotBlock) || errors.Is(err, user_model.ErrBlockOrganization) {
			ctx.Error(http.StatusBadRequest, "BlockUser", err)
		} else {
			ctx.Error(http.StatusInternalServerError, "BlockUser", err)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}

func UnblockUser(ctx *context.APIContext, doer, blocker *user_model.User) {
	blockee, err := user_model.GetUserByName(ctx, ctx.PathParam("username"))
	if err != nil {
		ctx.NotFound("GetUserByName", err)
		return
	}

	if err := user_service.UnblockUser(ctx, doer, blocker, blockee); err != nil {
		if errors.Is(err, user_model.ErrCanNotUnblock) || errors.Is(err, user_model.ErrBlockOrganization) {
			ctx.Error(http.StatusBadRequest, "UnblockUser", err)
		} else {
			ctx.Error(http.StatusInternalServerError, "UnblockUser", err)
		}
		return
	}

	ctx.Status(http.StatusNoContent)
}
