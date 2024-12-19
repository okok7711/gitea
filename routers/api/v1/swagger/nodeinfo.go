// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package swagger

import (
	api "github.com/okok7711/gitea/modules/structs"
)

// NodeInfo
// swagger:response NodeInfo
type swaggerResponseNodeInfo struct {
	// in:body
	Body api.NodeInfo `json:"body"`
}
