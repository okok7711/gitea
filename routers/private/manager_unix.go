// Copyright 2020 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build !windows

package private

import (
	"net/http"

	"github.com/okok7711/gitea/modules/graceful"
	"github.com/okok7711/gitea/services/context"
)

// Restart causes the server to perform a graceful restart
func Restart(ctx *context.PrivateContext) {
	graceful.GetManager().DoGracefulRestart()
	ctx.PlainText(http.StatusOK, "success")
}

// Shutdown causes the server to perform a graceful shutdown
func Shutdown(ctx *context.PrivateContext) {
	graceful.GetManager().DoGracefulShutdown()
	ctx.PlainText(http.StatusOK, "success")
}
