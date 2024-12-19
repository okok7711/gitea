// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"bytes"
	"io"
	"net/http"
	"path"

	"github.com/okok7711/gitea/models/renderhelper"
	"github.com/okok7711/gitea/modules/charset"
	"github.com/okok7711/gitea/modules/git"
	"github.com/okok7711/gitea/modules/log"
	"github.com/okok7711/gitea/modules/markup"
	"github.com/okok7711/gitea/modules/typesniffer"
	"github.com/okok7711/gitea/modules/util"
	"github.com/okok7711/gitea/services/context"
)

// RenderFile renders a file by repos path
func RenderFile(ctx *context.Context) {
	blob, err := ctx.Repo.Commit.GetBlobByPath(ctx.Repo.TreePath)
	if err != nil {
		if git.IsErrNotExist(err) {
			ctx.NotFound("GetBlobByPath", err)
		} else {
			ctx.ServerError("GetBlobByPath", err)
		}
		return
	}

	dataRc, err := blob.DataAsync()
	if err != nil {
		ctx.ServerError("DataAsync", err)
		return
	}
	defer dataRc.Close()

	buf := make([]byte, 1024)
	n, _ := util.ReadAtMost(dataRc, buf)
	buf = buf[:n]

	st := typesniffer.DetectContentType(buf)
	isTextFile := st.IsText()

	rd := charset.ToUTF8WithFallbackReader(io.MultiReader(bytes.NewReader(buf), dataRc), charset.ConvertOpts{})
	ctx.Resp.Header().Add("Content-Security-Policy", "frame-src 'self'; sandbox allow-scripts")

	if markupType := markup.DetectMarkupTypeByFileName(blob.Name()); markupType == "" {
		if isTextFile {
			_, _ = io.Copy(ctx.Resp, rd)
		} else {
			http.Error(ctx.Resp, "Unsupported file type render", http.StatusInternalServerError)
		}
		return
	}

	rctx := renderhelper.NewRenderContextRepoFile(ctx, ctx.Repo.Repository, renderhelper.RepoFileOptions{
		CurrentRefPath:  ctx.Repo.BranchNameSubURL(),
		CurrentTreePath: path.Dir(ctx.Repo.TreePath),
	}).WithRelativePath(ctx.Repo.TreePath).WithInStandalonePage(true)

	err = markup.Render(rctx, rd, ctx.Resp)
	if err != nil {
		log.Error("Failed to render file %q: %v", ctx.Repo.TreePath, err)
		http.Error(ctx.Resp, "Failed to render file", http.StatusInternalServerError)
		return
	}
}
