// Copyright 2017 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package misc

import (
	"net/http"

	"github.com/okok7711/gitea/modules/setting"
	"github.com/okok7711/gitea/modules/structs"
	"github.com/okok7711/gitea/services/context"
)

// Version shows the version of the Gitea server
func Version(ctx *context.APIContext) {
	// swagger:operation GET /version miscellaneous getVersion
	// ---
	// summary: Returns the version of the Gitea application
	// produces:
	// - application/json
	// responses:
	//   "200":
	//     "$ref": "#/responses/ServerVersion"
	ctx.JSON(http.StatusOK, &structs.ServerVersion{Version: setting.AppVer})
}
