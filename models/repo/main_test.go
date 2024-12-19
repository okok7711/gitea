// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo_test

import (
	"testing"

	"github.com/okok7711/gitea/models/unittest"

	_ "github.com/okok7711/gitea/models" // register table model
	_ "github.com/okok7711/gitea/models/actions"
	_ "github.com/okok7711/gitea/models/activities"
	_ "github.com/okok7711/gitea/models/perm/access" // register table model
	_ "github.com/okok7711/gitea/models/repo"        // register table model
	_ "github.com/okok7711/gitea/models/user"        // register table model
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
