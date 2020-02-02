// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

// TestIssueComment creat a issue and test comment creation/edit/deletion on it
func TestIssueComment(t *testing.T) {
	log.Println("== TestIssueComment ==")

	c := newTestClient()

	user, err := c.GetMyUserInfo()

	assert.NoError(t, err)
	repo, err := createTestRepo(t, "TestIssueCommentRepo", c)
	assert.NoError(t, err)
	issue1, err := c.CreateIssue(user.UserName, repo.Name, CreateIssueOption{Title: "issue1", Body: "body", Closed: false})
	assert.EqualValues(t, 1, issue1.Index)
	issue2, err := c.CreateIssue(user.UserName, repo.Name, CreateIssueOption{Title: "issue1", Body: "body", Closed: false})
	assert.EqualValues(t, 2, issue2.Index)
	assert.NoError(t, err)
	tUser2 := createTestUser(t, "Commenter2", c)
	tUser3 := createTestUser(t, "Commenter3", c)

	createOne := func(u *User, issue int64, text string) {
		c.sudo = u.UserName
		comment, e := c.CreateIssueComment(user.UserName, repo.Name, issue, CreateIssueCommentOption{Body: text})
		c.sudo = ""
		assert.NoError(t, e)
		assert.NotEmpty(t, comment)
		assert.EqualValues(t, text, comment.Body)
		assert.EqualValues(t, u.ID, comment.Poster.ID)
	}

	// CreateIssue
	createOne(user, 1, "what a nice issue")
	createOne(tUser2, 1, "dont think so")
	createOne(tUser3, 1, "weow weow")
	createOne(user, 1, "spam isn't it?")
	createOne(tUser3, 2, "hehe first commit")
	createOne(tUser2, 2, "second")
	createOne(user, 2, "3")

	assert.NoError(t, c.AdminDeleteUser(tUser3.UserName))

	// ListRepoIssueComments
	comments, err := c.ListRepoIssueComments(user.UserName, repo.Name)
	assert.NoError(t, err)
	assert.Len(t, comments, 7)

	// ListIssueComments
	comments, err = c.ListIssueComments(user.UserName, repo.Name, 2)
	assert.NoError(t, err)
	assert.Len(t, comments, 3)

	// GetIssueComment
	comment, err := c.GetIssueComment(user.UserName, repo.Name, comments[1].ID)
	assert.NoError(t, err)
	assert.EqualValues(t, comment.Poster.ID, comments[1].Poster.ID)
	assert.EqualValues(t, comment.Body, comments[1].Body)
	assert.EqualValues(t, comment.Updated.Unix(), comments[1].Updated.Unix())

	// EditIssueComment
	comment, err = c.EditIssueComment(user.UserName, repo.Name, comments[1].ID, EditIssueCommentOption{
		Body: "changed my mind",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, "changed my mind", comment.Body)

	// DeleteIssueComment
	assert.NoError(t, c.DeleteIssueComment(user.UserName, repo.Name, comments[1].ID))
	_, err = c.GetIssueComment(user.UserName, repo.Name, comments[1].ID)
	assert.Error(t, err)
}
