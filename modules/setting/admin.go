// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package setting

import (
	"github.com/okok7711/gitea/modules/container"
)

// Admin settings
var Admin struct {
	DisableRegularOrgCreation   bool
	DefaultEmailNotification    string
	UserDisabledFeatures        container.Set[string]
	ExternalUserDisableFeatures container.Set[string]
}

func loadAdminFrom(rootCfg ConfigProvider) {
	sec := rootCfg.Section("admin")
	Admin.DisableRegularOrgCreation = sec.Key("DISABLE_REGULAR_ORG_CREATION").MustBool(false)
	Admin.DefaultEmailNotification = sec.Key("DEFAULT_EMAIL_NOTIFICATIONS").MustString("enabled")
	Admin.UserDisabledFeatures = container.SetOf(sec.Key("USER_DISABLED_FEATURES").Strings(",")...)
	Admin.ExternalUserDisableFeatures = Admin.UserDisabledFeatures
}

const (
	UserFeatureDeletion          = "deletion"
	UserFeatureManageSSHKeys     = "manage_ssh_keys"
	UserFeatureManageGPGKeys     = "manage_gpg_keys"
	UserFeatureManageMFA         = "manage_mfa"
	UserFeatureManageCredentials = "manage_credentials"
	UserFeatureChangeUsername    = "change_username"
	UserFeatureChangeFullName    = "change_full_name"
)
