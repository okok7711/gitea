// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package system_test

import (
	"testing"

	"github.com/okok7711/gitea/models/unittest"

	_ "github.com/okok7711/gitea/models" // register models
	_ "github.com/okok7711/gitea/models/actions"
	_ "github.com/okok7711/gitea/models/activities"
	_ "github.com/okok7711/gitea/models/system" // register models of system
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
