// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"github.com/okok7711/gitea/modules/graceful"
	"github.com/okok7711/gitea/modules/log"
	"github.com/okok7711/gitea/modules/queue"
	"github.com/okok7711/gitea/modules/setting"
	notify_service "github.com/okok7711/gitea/services/notify"
)

func Init() {
	if !setting.Actions.Enabled {
		return
	}

	jobEmitterQueue = queue.CreateUniqueQueue(graceful.GetManager().ShutdownContext(), "actions_ready_job", jobEmitterQueueHandler)
	if jobEmitterQueue == nil {
		log.Fatal("Unable to create actions_ready_job queue")
	}
	go graceful.GetManager().RunWithCancel(jobEmitterQueue)

	notify_service.RegisterNotifier(NewNotifier())
}
