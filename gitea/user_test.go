// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyUser(t *testing.T) {
	log.Println("== TestMyUser ==")
	c := newTestClient()
	user, err := c.GetMyUserInfo()
	assert.NoError(t, err)

	assert.EqualValues(t, 1, user.ID)
	assert.EqualValues(t, "test01", user.UserName)
	assert.EqualValues(t, "test01@gitea.io", user.Email)
	assert.EqualValues(t, "", user.FullName)
	assert.EqualValues(t, getGiteaURL()+"/user/avatar/test01/-1", user.AvatarURL)
	assert.Equal(t, true, user.IsAdmin)
}

func TestUserApp(t *testing.T) {
	log.Println("== TestUserApp ==")
	c := newTestClient()

	result, err := c.ListAccessTokens(c.username, c.password)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.EqualValues(t, "gitea-admin", result[0].Name)

	t1, err := c.CreateAccessToken(c.username, c.password, CreateAccessTokenOption{Name: "TestCreateAccessToken"})
	assert.NoError(t, err)
	assert.EqualValues(t, "TestCreateAccessToken", t1.Name)
	result, _ = c.ListAccessTokens(c.username, c.password)
	assert.Len(t, result, 2)

	err = c.DeleteAccessToken(c.username, c.password, t1.ID)
	assert.NoError(t, err)
	result, _ = c.ListAccessTokens(c.username, c.password)
	assert.Len(t, result, 1)
}
