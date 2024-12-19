// Copyright 2024 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package markup

import (
	"bufio"
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/okok7711/gitea/models/perm/access"
	"github.com/okok7711/gitea/models/repo"
	"github.com/okok7711/gitea/models/unit"
	"github.com/okok7711/gitea/modules/charset"
	"github.com/okok7711/gitea/modules/gitrepo"
	"github.com/okok7711/gitea/modules/indexer/code"
	"github.com/okok7711/gitea/modules/markup"
	"github.com/okok7711/gitea/modules/setting"
	gitea_context "github.com/okok7711/gitea/services/context"
	"github.com/okok7711/gitea/services/repository/files"
)

func renderRepoFileCodePreview(ctx context.Context, opts markup.RenderCodePreviewOptions) (template.HTML, error) {
	opts.LineStop = max(opts.LineStop, opts.LineStart)
	lineCount := opts.LineStop - opts.LineStart + 1
	if lineCount <= 0 || lineCount > 140 /* GitHub at most show 140 lines */ {
		lineCount = 10
		opts.LineStop = opts.LineStart + lineCount
	}

	dbRepo, err := repo.GetRepositoryByOwnerAndName(ctx, opts.OwnerName, opts.RepoName)
	if err != nil {
		return "", err
	}

	webCtx, ok := ctx.Value(gitea_context.WebContextKey).(*gitea_context.Context)
	if !ok {
		return "", fmt.Errorf("context is not a web context")
	}
	doer := webCtx.Doer

	perms, err := access.GetUserRepoPermission(ctx, dbRepo, doer)
	if err != nil {
		return "", err
	}
	if !perms.CanRead(unit.TypeCode) {
		return "", fmt.Errorf("no permission")
	}

	gitRepo, err := gitrepo.OpenRepository(ctx, dbRepo)
	if err != nil {
		return "", err
	}
	defer gitRepo.Close()

	commit, err := gitRepo.GetCommit(opts.CommitID)
	if err != nil {
		return "", err
	}

	language, _ := files.TryGetContentLanguage(gitRepo, opts.CommitID, opts.FilePath)
	blob, err := commit.GetBlobByPath(opts.FilePath)
	if err != nil {
		return "", err
	}

	if blob.Size() > setting.UI.MaxDisplayFileSize {
		return "", fmt.Errorf("file is too large")
	}

	dataRc, err := blob.DataAsync()
	if err != nil {
		return "", err
	}
	defer dataRc.Close()

	reader := bufio.NewReader(dataRc)
	for i := 1; i < opts.LineStart; i++ {
		if _, err = reader.ReadBytes('\n'); err != nil {
			return "", err
		}
	}

	lineNums := make([]int, 0, lineCount)
	lineCodes := make([]string, 0, lineCount)
	for i := opts.LineStart; i <= opts.LineStop; i++ {
		line, err := reader.ReadString('\n')
		if err != nil && line == "" {
			break
		}

		lineNums = append(lineNums, i)
		lineCodes = append(lineCodes, line)
	}
	realLineStop := max(opts.LineStart, opts.LineStart+len(lineNums)-1)
	highlightLines := code.HighlightSearchResultCode(opts.FilePath, language, lineNums, strings.Join(lineCodes, ""))

	escapeStatus := &charset.EscapeStatus{}
	lineEscapeStatus := make([]*charset.EscapeStatus, len(highlightLines))
	for i, hl := range highlightLines {
		lineEscapeStatus[i], hl.FormattedContent = charset.EscapeControlHTML(hl.FormattedContent, webCtx.Base.Locale, charset.RuneNBSP)
		escapeStatus = escapeStatus.Or(lineEscapeStatus[i])
	}

	return webCtx.RenderToHTML("base/markup_codepreview", map[string]any{
		"FullURL":          opts.FullURL,
		"FilePath":         opts.FilePath,
		"LineStart":        opts.LineStart,
		"LineStop":         realLineStop,
		"RepoLink":         dbRepo.Link(),
		"CommitID":         opts.CommitID,
		"HighlightLines":   highlightLines,
		"EscapeStatus":     escapeStatus,
		"LineEscapeStatus": lineEscapeStatus,
	})
}
