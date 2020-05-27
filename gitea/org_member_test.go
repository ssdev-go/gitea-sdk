// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
// DeleteOrgMembership remove a member from an organization
func (c *Client) DeleteOrgMembership(org, user string) error {}


*/
func TestOrgMembership(t *testing.T) {
	log.Println("== TestOrgMembership ==")
	c := newTestClient()

	user := createTestUser(t, "org_mem_user", c)
	c.SetSudo(user.UserName)
	newOrg, err := c.CreateOrg(CreateOrgOption{UserName: "MemberOrg"})
	assert.NoError(t, err)
	assert.NotNil(t, newOrg)

	// Check func
	check, err := c.CheckPublicOrgMembership(newOrg.UserName, user.UserName)
	assert.NoError(t, err)
	assert.False(t, check)
	check, err = c.CheckOrgMembership(newOrg.UserName, user.UserName)
	assert.NoError(t, err)
	assert.True(t, check)

	err = c.SetPublicOrgMembership(newOrg.UserName, user.UserName, true)
	assert.NoError(t, err)
	check, err = c.CheckPublicOrgMembership(newOrg.UserName, user.UserName)
	assert.NoError(t, err)
	assert.True(t, check)

	u, err := c.ListOrgMembership(newOrg.UserName, ListOrgMembershipOption{})
	assert.NoError(t, err)
	assert.Len(t, u, 1)
	assert.EqualValues(t, user.UserName, u[0].UserName)
	u, err = c.ListPublicOrgMembership(newOrg.UserName, ListOrgMembershipOption{})
	assert.NoError(t, err)
	assert.Len(t, u, 1)
	assert.EqualValues(t, user.UserName, u[0].UserName)

	assert.Error(t, c.DeleteOrgMembership(newOrg.UserName, user.UserName))

	c.sudo = ""
	assert.Error(t, c.AdminDeleteUser(user.UserName))
	assert.NoError(t, c.DeleteOrg(newOrg.UserName))
	assert.NoError(t, c.AdminDeleteUser(user.UserName))
}
