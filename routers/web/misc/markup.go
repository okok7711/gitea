// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package misc

import (
	api "github.com/okok7711/gitea/modules/structs"
	"github.com/okok7711/gitea/modules/util"
	"github.com/okok7711/gitea/modules/web"
	"github.com/okok7711/gitea/routers/common"
	"github.com/okok7711/gitea/services/context"
)

// Markup render markup document to HTML
func Markup(ctx *context.Context) {
	form := web.GetForm(ctx).(*api.MarkupOption)
	mode := util.Iif(form.Wiki, "wiki", form.Mode) //nolint:staticcheck
	common.RenderMarkup(ctx.Base, ctx.Repo, mode, form.Text, form.Context, form.FilePath)
}
