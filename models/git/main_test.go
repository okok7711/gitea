// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package git_test

import (
	"testing"

	"github.com/okok7711/gitea/models/unittest"

	_ "github.com/okok7711/gitea/models"
	_ "github.com/okok7711/gitea/models/actions"
	_ "github.com/okok7711/gitea/models/activities"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
