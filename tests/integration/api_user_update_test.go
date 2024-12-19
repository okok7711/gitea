// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"net/http"
	"testing"

	auth_model "github.com/okok7711/gitea/models/auth"
	"github.com/okok7711/gitea/tests"
)

func TestAPIUpdateUser(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	normalUsername := "user2"
	session := loginUser(t, normalUsername)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteUser)

	req := NewRequestWithJSON(t, "PATCH", "/api/v1/user/settings", map[string]string{
		"website": "https://gitea.com",
	}).AddTokenAuth(token)
	MakeRequest(t, req, http.StatusOK)
}
