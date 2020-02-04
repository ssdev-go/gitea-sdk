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

func TestUserSearch(t *testing.T) {
	log.Println("== TestUserSearch ==")
	c := newTestClient()

	createTestUser(t, "tu1", c)
	createTestUser(t, "eatIt_2", c)
	createTestUser(t, "thirdIs3", c)
	createTestUser(t, "advancedUser", c)
	createTestUser(t, "1n2n3n", c)
	createTestUser(t, "otherIt", c)

	ul, err := c.SearchUsers(SearchUsersOption{KeyWord: "other"})
	assert.NoError(t, err)
	assert.Len(t, ul, 1)

	ul, err = c.SearchUsers(SearchUsersOption{KeyWord: "notInTESTcase"})
	assert.NoError(t, err)
	assert.Len(t, ul, 0)

	ul, err = c.SearchUsers(SearchUsersOption{KeyWord: "It"})
	assert.NoError(t, err)
	assert.Len(t, ul, 2)
}

func TestUserFollow(t *testing.T) {
	log.Println("== TestUserFollow ==")
	c := newTestClient()
	me, _ := c.GetMyUserInfo()

	uA := "uFollow_A"
	uB := "uFollow_B"
	uC := "uFollow_C"
	createTestUser(t, uA, c)
	createTestUser(t, uB, c)
	createTestUser(t, uC, c)

	// A follow ME
	// B follow C & ME
	// C follow A & B & ME
	c.sudo = uA
	err := c.Follow(me.UserName)
	assert.NoError(t, err)
	c.sudo = uB
	err = c.Follow(me.UserName)
	assert.NoError(t, err)
	err = c.Follow(uC)
	assert.NoError(t, err)
	c.sudo = uC
	err = c.Follow(me.UserName)
	assert.NoError(t, err)
	err = c.Follow(uA)
	assert.NoError(t, err)
	err = c.Follow(uB)
	assert.NoError(t, err)

	// C unfollow me
	err = c.Unfollow(me.UserName)
	assert.NoError(t, err)

	// ListMyFollowers of me
	c.sudo = ""
	f, err := c.ListMyFollowers(1)
	assert.NoError(t, err)
	assert.Len(t, f, 2)

	// ListFollowers of A
	f, err = c.ListFollowers(uA, 1)
	assert.NoError(t, err)
	assert.Len(t, f, 1)

	// ListMyFollowing of me
	f, err = c.ListMyFollowing(1)
	assert.NoError(t, err)
	assert.Len(t, f, 0)

	// ListFollowing of A
	f, err = c.ListFollowing(uA, 1)
	assert.NoError(t, err)
	assert.Len(t, f, 1)
	assert.EqualValues(t, me.ID, f[0].ID)

	assert.False(t, c.IsFollowing(uA))
	assert.True(t, c.IsUserFollowing(uB, uC))
}

func TestUserEmail(t *testing.T) {
	log.Println("== TestUserEmail ==")
	c := newTestClient()
	userN := "TestUserEmail"
	createTestUser(t, userN, c)
	c.sudo = userN

	// ListEmails
	el, err := c.ListEmails()
	assert.NoError(t, err)
	assert.Len(t, el, 1)
	assert.EqualValues(t, "testuseremail@gitea.io", el[0].Email)
	assert.True(t, el[0].Primary)

	// AddEmail
	mails := []string{"wow@mail.send", "speed@mail.me"}
	el, err = c.AddEmail(CreateEmailOption{Emails: mails})
	assert.NoError(t, err)
	assert.Len(t, el, 2)
	_, err = c.AddEmail(CreateEmailOption{Emails: []string{mails[1]}})
	assert.Error(t, err)
	el, err = c.ListEmails()
	assert.NoError(t, err)
	assert.Len(t, el, 3)

	// DeleteEmail
	err = c.DeleteEmail(DeleteEmailOption{Emails: []string{mails[1]}})
	assert.NoError(t, err)
	err = c.DeleteEmail(DeleteEmailOption{Emails: []string{"imaginary@e.de"}})
	assert.Error(t, err)

	el, err = c.ListEmails()
	assert.NoError(t, err)
	assert.Len(t, el, 2)
	err = c.DeleteEmail(DeleteEmailOption{Emails: []string{mails[0]}})
	assert.NoError(t, err)
}

func createTestUser(t *testing.T, username string, client *Client) *User {
	bFalse := false
	user, _ := client.GetUserInfo(username)
	if user.ID != 0 {
		return user
	}
	user, err := client.AdminCreateUser(CreateUserOption{Username: username, Password: username + "!1234", Email: username + "@gitea.io", MustChangePassword: &bFalse, SendNotify: bFalse})
	assert.NoError(t, err)
	return user
}
