// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"net/http"
	"testing"

	"github.com/okok7711/gitea/models/db"
	issues_model "github.com/okok7711/gitea/models/issues"
	repo_model "github.com/okok7711/gitea/models/repo"
	"github.com/okok7711/gitea/models/unittest"
	api "github.com/okok7711/gitea/modules/structs"
	"github.com/okok7711/gitea/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIPullCommits(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	pr := unittest.AssertExistsAndLoadBean(t, &issues_model.PullRequest{ID: 2})
	assert.NoError(t, pr.LoadIssue(db.DefaultContext))
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: pr.HeadRepoID})

	req := NewRequestf(t, http.MethodGet, "/api/v1/repos/%s/%s/pulls/%d/commits", repo.OwnerName, repo.Name, pr.Index)
	resp := MakeRequest(t, req, http.StatusOK)

	var commits []*api.Commit
	DecodeJSON(t, resp, &commits)

	if !assert.Len(t, commits, 2) {
		return
	}

	assert.Equal(t, "985f0301dba5e7b34be866819cd15ad3d8f508ee", commits[0].SHA)
	assert.Equal(t, "5c050d3b6d2db231ab1f64e324f1b6b9a0b181c2", commits[1].SHA)

	assert.NotEmpty(t, commits[0].Files)
	assert.NotEmpty(t, commits[1].Files)
	assert.NotNil(t, commits[0].RepoCommit.Verification)
	assert.NotNil(t, commits[1].RepoCommit.Verification)
}

// TODO add tests for already merged PR and closed PR
