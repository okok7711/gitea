// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package oauth2

import (
	"context"

	user_model "github.com/okok7711/gitea/models/user"
	"github.com/okok7711/gitea/services/auth/source/db"
)

// Authenticate falls back to the db authenticator
func (source *Source) Authenticate(ctx context.Context, user *user_model.User, login, password string) (*user_model.User, error) {
	return db.Authenticate(ctx, user, login, password)
}

// NB: Oauth2 does not implement LocalTwoFASkipper for password authentication
// as its password authentication drops to db authentication
