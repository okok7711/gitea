// Copyright 2016 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package routers

import (
	"context"
	"net/http"
	"reflect"
	"runtime"

	"github.com/okok7711/gitea/models"
	authmodel "github.com/okok7711/gitea/models/auth"
	"github.com/okok7711/gitea/modules/cache"
	"github.com/okok7711/gitea/modules/eventsource"
	"github.com/okok7711/gitea/modules/git"
	"github.com/okok7711/gitea/modules/highlight"
	"github.com/okok7711/gitea/modules/log"
	"github.com/okok7711/gitea/modules/markup"
	"github.com/okok7711/gitea/modules/markup/external"
	"github.com/okok7711/gitea/modules/setting"
	"github.com/okok7711/gitea/modules/ssh"
	"github.com/okok7711/gitea/modules/storage"
	"github.com/okok7711/gitea/modules/svg"
	"github.com/okok7711/gitea/modules/system"
	"github.com/okok7711/gitea/modules/templates"
	"github.com/okok7711/gitea/modules/translation"
	"github.com/okok7711/gitea/modules/util"
	"github.com/okok7711/gitea/modules/web"
	"github.com/okok7711/gitea/modules/web/routing"
	actions_router "github.com/okok7711/gitea/routers/api/actions"
	packages_router "github.com/okok7711/gitea/routers/api/packages"
	apiv1 "github.com/okok7711/gitea/routers/api/v1"
	"github.com/okok7711/gitea/routers/common"
	"github.com/okok7711/gitea/routers/private"
	web_routers "github.com/okok7711/gitea/routers/web"
	actions_service "github.com/okok7711/gitea/services/actions"
	asymkey_service "github.com/okok7711/gitea/services/asymkey"
	"github.com/okok7711/gitea/services/auth"
	"github.com/okok7711/gitea/services/auth/source/oauth2"
	"github.com/okok7711/gitea/services/automerge"
	"github.com/okok7711/gitea/services/cron"
	feed_service "github.com/okok7711/gitea/services/feed"
	indexer_service "github.com/okok7711/gitea/services/indexer"
	"github.com/okok7711/gitea/services/mailer"
	mailer_incoming "github.com/okok7711/gitea/services/mailer/incoming"
	markup_service "github.com/okok7711/gitea/services/markup"
	repo_migrations "github.com/okok7711/gitea/services/migrations"
	mirror_service "github.com/okok7711/gitea/services/mirror"
	"github.com/okok7711/gitea/services/oauth2_provider"
	pull_service "github.com/okok7711/gitea/services/pull"
	release_service "github.com/okok7711/gitea/services/release"
	repo_service "github.com/okok7711/gitea/services/repository"
	"github.com/okok7711/gitea/services/repository/archiver"
	"github.com/okok7711/gitea/services/task"
	"github.com/okok7711/gitea/services/uinotification"
	"github.com/okok7711/gitea/services/webhook"
)

func mustInit(fn func() error) {
	err := fn()
	if err != nil {
		ptr := reflect.ValueOf(fn).Pointer()
		fi := runtime.FuncForPC(ptr)
		log.Fatal("%s failed: %v", fi.Name(), err)
	}
}

func mustInitCtx(ctx context.Context, fn func(ctx context.Context) error) {
	err := fn(ctx)
	if err != nil {
		ptr := reflect.ValueOf(fn).Pointer()
		fi := runtime.FuncForPC(ptr)
		log.Fatal("%s(ctx) failed: %v", fi.Name(), err)
	}
}

func syncAppConfForGit(ctx context.Context) error {
	runtimeState := new(system.RuntimeState)
	if err := system.AppState.Get(ctx, runtimeState); err != nil {
		return err
	}

	updated := false
	if runtimeState.LastAppPath != setting.AppPath {
		log.Info("AppPath changed from '%s' to '%s'", runtimeState.LastAppPath, setting.AppPath)
		runtimeState.LastAppPath = setting.AppPath
		updated = true
	}
	if runtimeState.LastCustomConf != setting.CustomConf {
		log.Info("CustomConf changed from '%s' to '%s'", runtimeState.LastCustomConf, setting.CustomConf)
		runtimeState.LastCustomConf = setting.CustomConf
		updated = true
	}

	if updated {
		log.Info("re-sync repository hooks ...")
		mustInitCtx(ctx, repo_service.SyncRepositoryHooks)

		log.Info("re-write ssh public keys ...")
		mustInitCtx(ctx, asymkey_service.RewriteAllPublicKeys)

		return system.AppState.Set(ctx, runtimeState)
	}
	return nil
}

func InitWebInstallPage(ctx context.Context) {
	translation.InitLocales(ctx)
	setting.LoadSettingsForInstall()
	mustInit(svg.Init)
}

// InitWebInstalled is for global installed configuration.
func InitWebInstalled(ctx context.Context) {
	mustInitCtx(ctx, git.InitFull)
	log.Info("Git version: %s (home: %s)", git.DefaultFeatures().VersionInfo(), git.HomeDir())
	if !git.DefaultFeatures().SupportHashSha256 {
		log.Warn("sha256 hash support is disabled - requires Git >= 2.42." + util.Iif(git.DefaultFeatures().UsingGogit, " Gogit is currently unsupported.", ""))
	}

	// Setup i18n
	translation.InitLocales(ctx)

	setting.LoadSettings()
	mustInit(storage.Init)

	mailer.NewContext(ctx)
	mustInit(cache.Init)
	mustInit(feed_service.Init)
	mustInit(uinotification.Init)
	mustInitCtx(ctx, archiver.Init)

	highlight.NewContext()
	external.RegisterRenderers()
	markup.Init(markup_service.ProcessorHelper())

	if setting.EnableSQLite3 {
		log.Info("SQLite3 support is enabled")
	} else if setting.Database.Type.IsSQLite3() {
		log.Fatal("SQLite3 support is disabled, but it is used for database setting. Please get or build a Gitea release with SQLite3 support.")
	}

	mustInitCtx(ctx, common.InitDBEngine)
	log.Info("ORM engine initialization successful!")
	mustInit(system.Init)
	mustInitCtx(ctx, oauth2.Init)
	mustInitCtx(ctx, oauth2_provider.Init)
	mustInit(release_service.Init)

	mustInitCtx(ctx, models.Init)
	mustInitCtx(ctx, authmodel.Init)
	mustInitCtx(ctx, repo_service.Init)

	// Booting long running goroutines.
	mustInit(indexer_service.Init)

	mirror_service.InitSyncMirrors()
	mustInit(webhook.Init)
	mustInit(pull_service.Init)
	mustInit(automerge.Init)
	mustInit(task.Init)
	mustInit(repo_migrations.Init)
	eventsource.GetManager().Init()
	mustInitCtx(ctx, mailer_incoming.Init)

	mustInitCtx(ctx, syncAppConfForGit)

	mustInit(ssh.Init)

	auth.Init()
	mustInit(svg.Init)

	actions_service.Init()

	mustInit(repo_service.InitLicenseClassifier)

	// Finally start up the cron
	cron.NewContext(ctx)
}

// NormalRoutes represents non install routes
func NormalRoutes() *web.Router {
	_ = templates.HTMLRenderer()
	r := web.NewRouter()
	r.Use(common.ProtocolMiddlewares()...)

	r.Mount("/", web_routers.Routes())
	r.Mount("/api/v1", apiv1.Routes())
	r.Mount("/api/internal", private.Routes())

	r.Post("/-/fetch-redirect", common.FetchRedirectDelegate)

	if setting.Packages.Enabled {
		// This implements package support for most package managers
		r.Mount("/api/packages", packages_router.CommonRoutes())
		// This implements the OCI API, this container registry "/v2" endpoint must be in the root of the site.
		// If site admin deploys Gitea in a sub-path, they must configure their reverse proxy to map the "https://host/v2" endpoint to Gitea.
		r.Mount("/v2", packages_router.ContainerRoutes())
	}

	if setting.Actions.Enabled {
		prefix := "/api/actions"
		r.Mount(prefix, actions_router.Routes(prefix))

		// TODO: Pipeline api used for runner internal communication with gitea server. but only artifact is used for now.
		// In Github, it uses ACTIONS_RUNTIME_URL=https://pipelines.actions.githubusercontent.com/fLgcSHkPGySXeIFrg8W8OBSfeg3b5Fls1A1CwX566g8PayEGlg/
		// TODO: this prefix should be generated with a token string with runner ?
		prefix = "/api/actions_pipeline"
		r.Mount(prefix, actions_router.ArtifactsRoutes(prefix))
		prefix = actions_router.ArtifactV4RouteBase
		r.Mount(prefix, actions_router.ArtifactsV4Routes(prefix))
	}

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		routing.UpdateFuncInfo(req.Context(), routing.GetFuncInfo(http.NotFound, "GlobalNotFound"))
		http.NotFound(w, req)
	})
	return r
}
