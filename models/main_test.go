// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package models

import (
	"testing"

	activities_model "github.com/okok7711/gitea/models/activities"
	"github.com/okok7711/gitea/models/organization"
	repo_model "github.com/okok7711/gitea/models/repo"
	"github.com/okok7711/gitea/models/unittest"
	user_model "github.com/okok7711/gitea/models/user"

	_ "github.com/okok7711/gitea/models/actions"
	_ "github.com/okok7711/gitea/models/system"

	"github.com/stretchr/testify/assert"
)

// TestFixturesAreConsistent assert that test fixtures are consistent
func TestFixturesAreConsistent(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())
	unittest.CheckConsistencyFor(t,
		&user_model.User{},
		&repo_model.Repository{},
		&organization.Team{},
		&activities_model.Action{})
}

func TestMain(m *testing.M) {
	unittest.MainTest(m)
}
