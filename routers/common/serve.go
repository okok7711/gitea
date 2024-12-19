// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package common

import (
	"io"
	"time"

	"github.com/okok7711/gitea/modules/git"
	"github.com/okok7711/gitea/modules/httpcache"
	"github.com/okok7711/gitea/modules/httplib"
	"github.com/okok7711/gitea/modules/log"
	"github.com/okok7711/gitea/services/context"
)

// ServeBlob download a git.Blob
func ServeBlob(ctx *context.Base, filePath string, blob *git.Blob, lastModified *time.Time) error {
	if httpcache.HandleGenericETagTimeCache(ctx.Req, ctx.Resp, `"`+blob.ID.String()+`"`, lastModified) {
		return nil
	}

	dataRc, err := blob.DataAsync()
	if err != nil {
		return err
	}
	defer func() {
		if err = dataRc.Close(); err != nil {
			log.Error("ServeBlob: Close: %v", err)
		}
	}()

	httplib.ServeContentByReader(ctx.Req, ctx.Resp, filePath, blob.Size(), dataRc)
	return nil
}

func ServeContentByReader(ctx *context.Base, filePath string, size int64, reader io.Reader) {
	httplib.ServeContentByReader(ctx.Req, ctx.Resp, filePath, size, reader)
}

func ServeContentByReadSeeker(ctx *context.Base, filePath string, modTime *time.Time, reader io.ReadSeeker) {
	httplib.ServeContentByReadSeeker(ctx.Req, ctx.Resp, filePath, modTime, reader)
}
