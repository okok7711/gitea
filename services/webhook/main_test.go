// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package webhook

import (
	"testing"

	"github.com/okok7711/gitea/models/unittest"
	"github.com/okok7711/gitea/modules/hostmatcher"
	"github.com/okok7711/gitea/modules/setting"

	_ "github.com/okok7711/gitea/models"
	_ "github.com/okok7711/gitea/models/actions"
)

func TestMain(m *testing.M) {
	// for tests, allow only loopback IPs
	setting.Webhook.AllowedHostList = hostmatcher.MatchBuiltinLoopback
	unittest.MainTest(m, &unittest.TestOptions{
		SetUp: func() error {
			setting.LoadQueueSettings()
			return Init()
		},
	})
}
