// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package issues_test

import (
	"testing"

	"github.com/okok7711/gitea/models/db"
	issues_model "github.com/okok7711/gitea/models/issues"
	"github.com/okok7711/gitea/models/unittest"
	user_model "github.com/okok7711/gitea/models/user"

	"github.com/stretchr/testify/assert"
)

func TestNewIssueLabelsScope(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 18})
	label1 := unittest.AssertExistsAndLoadBean(t, &issues_model.Label{ID: 7})
	label2 := unittest.AssertExistsAndLoadBean(t, &issues_model.Label{ID: 8})
	doer := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	assert.NoError(t, issues_model.NewIssueLabels(db.DefaultContext, issue, []*issues_model.Label{label1, label2}, doer))

	assert.Len(t, issue.Labels, 1)
	assert.Equal(t, label2.ID, issue.Labels[0].ID)
}
