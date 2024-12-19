// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package avatars_test

import (
	"testing"

	"github.com/okok7711/gitea/models/unittest"

	_ "github.com/okok7711/gitea/models"
	_ "github.com/okok7711/gitea/models/activities"
	_ "github.com/okok7711/gitea/models/perm/access"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
