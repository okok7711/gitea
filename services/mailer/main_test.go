// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package mailer

import (
	"testing"

	"github.com/okok7711/gitea/models/unittest"

	_ "github.com/okok7711/gitea/models"
	_ "github.com/okok7711/gitea/models/actions"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
