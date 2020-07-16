// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminOrg(t *testing.T) {
	log.Println("== TestAdminOrg ==")
	c := newTestClient()
	user, err := c.GetMyUserInfo()
	assert.NoError(t, err)

	orgName := "NewTestOrg"
	newOrg, err := c.AdminCreateOrg(user.UserName, CreateOrgOption{
		UserName:    orgName,
		FullName:    orgName + " FullName",
		Description: "test adminCreateOrg",
		Visibility:  VisibleTypePublic,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, newOrg)
	assert.EqualValues(t, orgName, newOrg.UserName)

	orgs, err := c.AdminListOrgs(AdminListOrgsOptions{})
	assert.NoError(t, err)
	assert.Len(t, orgs, 1)
	assert.EqualValues(t, newOrg.ID, orgs[0].ID)

	err = c.DeleteOrg(orgName)
	assert.NoError(t, err)
}
