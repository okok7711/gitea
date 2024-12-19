// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package user

import (
	"encoding/base64"
	"net/http"

	api "github.com/okok7711/gitea/modules/structs"
	"github.com/okok7711/gitea/modules/web"
	"github.com/okok7711/gitea/services/context"
	user_service "github.com/okok7711/gitea/services/user"
)

// UpdateAvatar updates the Avatar of an User
func UpdateAvatar(ctx *context.APIContext) {
	// swagger:operation POST /user/avatar user userUpdateAvatar
	// ---
	// summary: Update Avatar
	// produces:
	// - application/json
	// parameters:
	// - name: body
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/UpdateUserAvatarOption"
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	form := web.GetForm(ctx).(*api.UpdateUserAvatarOption)

	content, err := base64.StdEncoding.DecodeString(form.Image)
	if err != nil {
		ctx.Error(http.StatusBadRequest, "DecodeImage", err)
		return
	}

	err = user_service.UploadAvatar(ctx, ctx.Doer, content)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "UploadAvatar", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// DeleteAvatar deletes the Avatar of an User
func DeleteAvatar(ctx *context.APIContext) {
	// swagger:operation DELETE /user/avatar user userDeleteAvatar
	// ---
	// summary: Delete Avatar
	// produces:
	// - application/json
	// responses:
	//   "204":
	//     "$ref": "#/responses/empty"
	err := user_service.DeleteAvatar(ctx, ctx.Doer)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "DeleteAvatar", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
