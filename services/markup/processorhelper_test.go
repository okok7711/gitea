// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package markup

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/okok7711/gitea/models/db"
	"github.com/okok7711/gitea/models/unittest"
	"github.com/okok7711/gitea/models/user"
	gitea_context "github.com/okok7711/gitea/services/context"
	"github.com/okok7711/gitea/services/contexttest"

	"github.com/stretchr/testify/assert"
)

func TestProcessorHelper(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	userPublic := "user1"
	userPrivate := "user31"
	userLimited := "user33"
	userNoSuch := "no-such-user"

	unittest.AssertCount(t, &user.User{Name: userPublic}, 1)
	unittest.AssertCount(t, &user.User{Name: userPrivate}, 1)
	unittest.AssertCount(t, &user.User{Name: userLimited}, 1)
	unittest.AssertCount(t, &user.User{Name: userNoSuch}, 0)

	// when using general context, use user's visibility to check
	assert.True(t, ProcessorHelper().IsUsernameMentionable(context.Background(), userPublic))
	assert.False(t, ProcessorHelper().IsUsernameMentionable(context.Background(), userLimited))
	assert.False(t, ProcessorHelper().IsUsernameMentionable(context.Background(), userPrivate))
	assert.False(t, ProcessorHelper().IsUsernameMentionable(context.Background(), userNoSuch))

	// when using web context, use user.IsUserVisibleToViewer to check
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	base, baseCleanUp := gitea_context.NewBaseContext(httptest.NewRecorder(), req)
	defer baseCleanUp()
	giteaCtx := gitea_context.NewWebContext(base, &contexttest.MockRender{}, nil)

	assert.True(t, ProcessorHelper().IsUsernameMentionable(giteaCtx, userPublic))
	assert.False(t, ProcessorHelper().IsUsernameMentionable(giteaCtx, userPrivate))

	giteaCtx.Doer, err = user.GetUserByName(db.DefaultContext, userPrivate)
	assert.NoError(t, err)
	assert.True(t, ProcessorHelper().IsUsernameMentionable(giteaCtx, userPublic))
	assert.True(t, ProcessorHelper().IsUsernameMentionable(giteaCtx, userPrivate))
}
