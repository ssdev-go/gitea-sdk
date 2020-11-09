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
	user, _, err := c.GetMyUserInfo()
	assert.NoError(t, err)

	orgName := "NewTestOrg"
	newOrg, _, err := c.AdminCreateOrg(user.UserName, CreateOrgOption{
		Name:        orgName,
		FullName:    orgName + " FullName",
		Description: "test adminCreateOrg",
		Visibility:  VisibleTypePublic,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, newOrg)
	assert.EqualValues(t, orgName, newOrg.UserName)

	orgs, _, err := c.AdminListOrgs(AdminListOrgsOptions{})
	assert.NoError(t, err)
	assert.Len(t, orgs, 1)
	assert.EqualValues(t, newOrg.ID, orgs[0].ID)

	_, err = c.DeleteOrg(orgName)
	assert.NoError(t, err)
}

func TestAdminCronTasks(t *testing.T) {
	log.Println("== TestAdminCronTasks ==")
	c := newTestClient()

	tasks, _, err := c.ListCronTasks(ListCronTaskOptions{})
	assert.NoError(t, err)
	assert.Len(t, tasks, 16)
	_, err = c.RunCronTasks(tasks[0].Name)
	assert.NoError(t, err)
}
