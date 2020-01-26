// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepo(t *testing.T) {
	log.Println("== TestCreateRepo ==")
	c := newTestClient()
	user, err := c.GetMyUserInfo()
	assert.NoError(t, err)

	var repoName = "test1"
	_, err = c.GetRepo(user.UserName, repoName)
	if err != nil {
		repo, err := c.CreateRepo(CreateRepoOption{
			Name: repoName,
		})
		assert.NoError(t, err)
		assert.NotNil(t, repo)
	}

	err = c.DeleteRepo(user.UserName, repoName)
	assert.NoError(t, err)
}

func TestDeleteRepo(t *testing.T) {
	log.Println("== TestDeleteRepo ==")
	c := newTestClient()
	repo, _ := createTestRepo(t, "TestDeleteRepo", c)
	assert.NoError(t, c.DeleteRepo(repo.Owner.UserName, repo.Name))
}

// standard func to create a init repo for test routines
func createTestRepo(t *testing.T, name string, c *Client) (*Repository, error) {
	user, uErr := c.GetMyUserInfo()
	assert.NoError(t, uErr)
	_, err := c.GetRepo(user.UserName, name)
	if err == nil {
		_ = c.DeleteRepo(user.UserName, name)
	}
	repo, err := c.CreateRepo(CreateRepoOption{
		Name:        name,
		Description: "A test Repo: " + name,
		AutoInit:    true,
		Gitignores:  "C,C++",
		License:     "MIT",
		Readme:      "Default",
		IssueLabels: "Default",
		Private:     false,
	})
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	return repo, err
}
