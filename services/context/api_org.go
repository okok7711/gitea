// Copyright 2016 The Gogs Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package context

import "github.com/okok7711/gitea/models/organization"

// APIOrganization contains organization and team
type APIOrganization struct {
	Organization *organization.Organization
	Team         *organization.Team
}
