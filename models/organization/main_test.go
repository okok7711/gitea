// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package organization_test

import (
	"testing"

	"github.com/okok7711/gitea/models/unittest"

	_ "github.com/okok7711/gitea/models"
	_ "github.com/okok7711/gitea/models/actions"
	_ "github.com/okok7711/gitea/models/activities"
	_ "github.com/okok7711/gitea/models/organization"
	_ "github.com/okok7711/gitea/models/repo"
	_ "github.com/okok7711/gitea/models/user"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
