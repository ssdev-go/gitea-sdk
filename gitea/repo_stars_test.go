// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestRepoStaring(t *testing.T) {
	log.Println("== TestRepoStaring ==")

	// init user2
	c := newTestClient()

	user1, _, err := c.GetMyUserInfo()
	assert.NoError(t, err)

	userA := createTestUser(t, "stargazer_a", c)
	userB := createTestUser(t, "stargazer_b", c)

	repo, _ := createTestRepo(t, "toStar", c)
	if repo == nil {
		t.Skip()
	}

	is, _, err := c.IsRepoStarring(repo.Owner.UserName, repo.Name)
	assert.NoError(t, err)
	assert.False(t, is)

	repos, _, err := c.GetMyStarredRepos()
	assert.NoError(t, err)
	assert.Len(t, repos, 0)

	_, err = c.StarRepo(repo.Owner.UserName, repo.Name)
	assert.NoError(t, err)
	c.SetSudo(userA.UserName)
	_, err = c.StarRepo(repo.Owner.UserName, repo.Name)
	assert.NoError(t, err)
	c.SetSudo(userB.UserName)
	_, err = c.StarRepo(repo.Owner.UserName, repo.Name)
	assert.NoError(t, err)

	users, _, err := c.ListRepoStargazers(repo.Owner.UserName, repo.Name, ListStargazersOptions{})
	assert.NoError(t, err)
	assert.Len(t, users, 3)
	assert.EqualValues(t, user1.UserName, users[0].UserName)

	_, err = c.UnStarRepo(repo.Owner.UserName, repo.Name)
	assert.NoError(t, err)
	_, err = c.UnStarRepo(repo.Owner.UserName, repo.Name)
	assert.NoError(t, err)

	c.SetSudo("")

	users, _, err = c.ListRepoStargazers(repo.Owner.UserName, repo.Name, ListStargazersOptions{})
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	repos, _, err = c.GetMyStarredRepos()
	assert.NoError(t, err)
	assert.Len(t, repos, 1)

	reposNew, _, err := c.GetStarredRepos(user1.UserName)
	assert.NoError(t, err)
	assert.Len(t, repos, 1)
	assert.EqualValues(t, repos, reposNew)
}
