// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package cmd

import (
	"github.com/okok7711/gitea/modules/graceful"
	asymkey_service "github.com/okok7711/gitea/services/asymkey"
	repo_service "github.com/okok7711/gitea/services/repository"

	"github.com/urfave/cli/v2"
)

var (
	microcmdRegenHooks = &cli.Command{
		Name:   "hooks",
		Usage:  "Regenerate git-hooks",
		Action: runRegenerateHooks,
	}

	microcmdRegenKeys = &cli.Command{
		Name:   "keys",
		Usage:  "Regenerate authorized_keys file",
		Action: runRegenerateKeys,
	}
)

func runRegenerateHooks(_ *cli.Context) error {
	ctx, cancel := installSignals()
	defer cancel()

	if err := initDB(ctx); err != nil {
		return err
	}
	return repo_service.SyncRepositoryHooks(graceful.GetManager().ShutdownContext())
}

func runRegenerateKeys(_ *cli.Context) error {
	ctx, cancel := installSignals()
	defer cancel()

	if err := initDB(ctx); err != nil {
		return err
	}
	return asymkey_service.RewriteAllPublicKeys(ctx)
}
