// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package models

import (
	"testing"

	"github.com/okok7711/gitea/models/db"
	"github.com/okok7711/gitea/models/unittest"

	"github.com/stretchr/testify/assert"
)

func TestCheckRepoStats(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	assert.NoError(t, CheckRepoStats(db.DefaultContext))
}

func TestDoctorUserStarNum(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	assert.NoError(t, DoctorUserStarNum(db.DefaultContext))
}
