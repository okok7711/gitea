// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package org

import (
	"net/http"

	"github.com/okok7711/gitea/modules/base"
	shared_user "github.com/okok7711/gitea/routers/web/shared/user"
	"github.com/okok7711/gitea/services/context"
)

const (
	tplSettingsBlockedUsers base.TplName = "org/settings/blocked_users"
)

func BlockedUsers(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("user.block.list")
	ctx.Data["PageIsOrgSettings"] = true
	ctx.Data["PageIsSettingsBlockedUsers"] = true

	shared_user.BlockedUsers(ctx, ctx.ContextUser)
	if ctx.Written() {
		return
	}

	ctx.HTML(http.StatusOK, tplSettingsBlockedUsers)
}

func BlockedUsersPost(ctx *context.Context) {
	shared_user.BlockedUsersPost(ctx, ctx.ContextUser)
	if ctx.Written() {
		return
	}

	ctx.Redirect(ctx.ContextUser.OrganisationLink() + "/settings/blocked_users")
}
