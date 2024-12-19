// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package db_test

import (
	"testing"

	"github.com/okok7711/gitea/models/unittest"

	_ "github.com/okok7711/gitea/models"
	_ "github.com/okok7711/gitea/models/repo"
)

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
