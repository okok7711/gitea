// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package admin

import (
	"github.com/okok7711/gitea/modules/setting"
	"github.com/okok7711/gitea/services/context"
)

func RedirectToDefaultSetting(ctx *context.Context) {
	ctx.Redirect(setting.AppSubURL + "/-/admin/actions/runners")
}
