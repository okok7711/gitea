// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"io"
	"net/http"
	"testing"

	"github.com/okok7711/gitea/modules/setting"
	"github.com/okok7711/gitea/modules/test"
	"github.com/okok7711/gitea/routers"
	"github.com/okok7711/gitea/routers/web"
	"github.com/okok7711/gitea/tests"

	"github.com/stretchr/testify/assert"
)

func TestRepoDownloadArchive(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	defer test.MockVariableValue(&setting.EnableGzip, true)()
	defer test.MockVariableValue(&web.GzipMinSize, 10)()
	defer test.MockVariableValue(&testWebRoutes, routers.NormalRoutes())()

	req := NewRequest(t, "GET", "/user2/repo1/archive/master.zip")
	req.Header.Set("Accept-Encoding", "gzip")
	resp := MakeRequest(t, req, http.StatusOK)
	bs, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Empty(t, resp.Header().Get("Content-Encoding"))
	assert.Len(t, bs, 320)
}
