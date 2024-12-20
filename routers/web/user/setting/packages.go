// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import (
	"net/http"
	"strings"

	user_model "github.com/okok7711/gitea/models/user"
	"github.com/okok7711/gitea/modules/base"
	chef_module "github.com/okok7711/gitea/modules/packages/chef"
	"github.com/okok7711/gitea/modules/setting"
	"github.com/okok7711/gitea/modules/util"
	shared "github.com/okok7711/gitea/routers/web/shared/packages"
	"github.com/okok7711/gitea/services/context"
)

const (
	tplSettingsPackages            base.TplName = "user/settings/packages"
	tplSettingsPackagesRuleEdit    base.TplName = "user/settings/packages_cleanup_rules_edit"
	tplSettingsPackagesRulePreview base.TplName = "user/settings/packages_cleanup_rules_preview"
)

func Packages(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.SetPackagesContext(ctx, ctx.Doer)

	ctx.HTML(http.StatusOK, tplSettingsPackages)
}

func PackagesRuleAdd(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.SetRuleAddContext(ctx)

	ctx.HTML(http.StatusOK, tplSettingsPackagesRuleEdit)
}

func PackagesRuleEdit(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.SetRuleEditContext(ctx, ctx.Doer)

	ctx.HTML(http.StatusOK, tplSettingsPackagesRuleEdit)
}

func PackagesRuleAddPost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("settings")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.PerformRuleAddPost(
		ctx,
		ctx.Doer,
		setting.AppSubURL+"/user/settings/packages",
		tplSettingsPackagesRuleEdit,
	)
}

func PackagesRuleEditPost(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.PerformRuleEditPost(
		ctx,
		ctx.Doer,
		setting.AppSubURL+"/user/settings/packages",
		tplSettingsPackagesRuleEdit,
	)
}

func PackagesRulePreview(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true
	ctx.Data["UserDisabledFeatures"] = user_model.DisabledFeaturesWithLoginType(ctx.Doer)

	shared.SetRulePreviewContext(ctx, ctx.Doer)

	ctx.HTML(http.StatusOK, tplSettingsPackagesRulePreview)
}

func InitializeCargoIndex(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true

	shared.InitializeCargoIndex(ctx, ctx.Doer)

	ctx.Redirect(setting.AppSubURL + "/user/settings/packages")
}

func RebuildCargoIndex(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("packages.title")
	ctx.Data["PageIsSettingsPackages"] = true

	shared.RebuildCargoIndex(ctx, ctx.Doer)

	ctx.Redirect(setting.AppSubURL + "/user/settings/packages")
}

func RegenerateChefKeyPair(ctx *context.Context) {
	priv, pub, err := util.GenerateKeyPair(chef_module.KeyBits)
	if err != nil {
		ctx.ServerError("GenerateKeyPair", err)
		return
	}

	if err := user_model.SetUserSetting(ctx, ctx.Doer.ID, chef_module.SettingPublicPem, pub); err != nil {
		ctx.ServerError("SetUserSetting", err)
		return
	}

	ctx.ServeContent(strings.NewReader(priv), &context.ServeHeaderOptions{
		ContentType: "application/x-pem-file",
		Filename:    ctx.Doer.Name + ".priv",
	})
}
