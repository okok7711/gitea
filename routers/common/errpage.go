// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package common

import (
	"fmt"
	"net/http"

	user_model "github.com/okok7711/gitea/models/user"
	"github.com/okok7711/gitea/modules/base"
	"github.com/okok7711/gitea/modules/httpcache"
	"github.com/okok7711/gitea/modules/log"
	"github.com/okok7711/gitea/modules/setting"
	"github.com/okok7711/gitea/modules/templates"
	"github.com/okok7711/gitea/modules/web/middleware"
	"github.com/okok7711/gitea/modules/web/routing"
	"github.com/okok7711/gitea/services/context"
)

const tplStatus500 base.TplName = "status/500"

// RenderPanicErrorPage renders a 500 page, and it never panics
func RenderPanicErrorPage(w http.ResponseWriter, req *http.Request, err any) {
	combinedErr := fmt.Sprintf("%v\n%s", err, log.Stack(2))
	log.Error("PANIC: %s", combinedErr)

	defer func() {
		if err := recover(); err != nil {
			log.Error("Panic occurs again when rendering error page: %v. Stack:\n%s", err, log.Stack(2))
		}
	}()

	routing.UpdatePanicError(req.Context(), err)

	httpcache.SetCacheControlInHeader(w.Header(), 0, "no-transform")
	w.Header().Set(`X-Frame-Options`, setting.CORSConfig.XFrameOptions)

	tmplCtx := context.TemplateContext{}
	tmplCtx["Locale"] = middleware.Locale(w, req)
	ctxData := middleware.GetContextData(req.Context())

	// This recovery handler could be called without Gitea's web context, so we shouldn't touch that context too much.
	// Otherwise, the 500-page may cause new panics, eg: cache.GetContextWithData, it makes the developer&users couldn't find the original panic.
	user, _ := ctxData[middleware.ContextDataKeySignedUser].(*user_model.User)
	if !setting.IsProd || (user != nil && user.IsAdmin) {
		ctxData["ErrorMsg"] = "PANIC: " + combinedErr
	}

	err = templates.HTMLRenderer().HTML(w, http.StatusInternalServerError, string(tplStatus500), ctxData, tmplCtx)
	if err != nil {
		log.Error("Error occurs again when rendering error page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal server error, please collect error logs and report to Gitea issue tracker"))
	}
}
