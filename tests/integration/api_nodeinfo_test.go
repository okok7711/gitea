// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/okok7711/gitea/modules/setting"
	api "github.com/okok7711/gitea/modules/structs"
	"github.com/okok7711/gitea/routers"

	"github.com/stretchr/testify/assert"
)

func TestNodeinfo(t *testing.T) {
	setting.Federation.Enabled = true
	testWebRoutes = routers.NormalRoutes()
	defer func() {
		setting.Federation.Enabled = false
		testWebRoutes = routers.NormalRoutes()
	}()

	onGiteaRun(t, func(*testing.T, *url.URL) {
		req := NewRequest(t, "GET", "/api/v1/nodeinfo")
		resp := MakeRequest(t, req, http.StatusOK)
		VerifyJSONSchema(t, resp, "nodeinfo_2.1.json")

		var nodeinfo api.NodeInfo
		DecodeJSON(t, resp, &nodeinfo)
		assert.True(t, nodeinfo.OpenRegistrations)
		assert.Equal(t, "gitea", nodeinfo.Software.Name)
		assert.Equal(t, 29, nodeinfo.Usage.Users.Total)
		assert.Equal(t, 22, nodeinfo.Usage.LocalPosts)
		assert.Equal(t, 3, nodeinfo.Usage.LocalComments)
	})
}
