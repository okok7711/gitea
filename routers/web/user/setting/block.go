// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import (
	"net/http"

	user_model "github.com/okok7711/gitea/models/user"
	"github.com/okok7711/gitea/modules/base"
	"github.com/okok7711/gitea/modules/setting"
	shared_user "github.com/okok7711/gitea/routers/web/shared/user"
	"github.com/okok7711/gitea/services/context"
)

const (
	tplSettingsBlockedUsers base.TplName = "user/settings/blocked_users"
)

func BlockedUsers(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("user.block.list")
	ctx.Data["PageIsSettingsBlockedUsers"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared_user.BlockedUsers(ctx, ctx.Doer)
	if ctx.Written() {
		return
	}

	ctx.HTML(http.StatusOK, tplSettingsBlockedUsers)
}

func BlockedUsersPost(ctx *context.Context) {
	shared_user.BlockedUsersPost(ctx, ctx.Doer)
	if ctx.Written() {
		return
	}

	ctx.Redirect(setting.AppSubURL + "/user/settings/blocked_users")
}
