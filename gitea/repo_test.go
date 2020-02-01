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

func TestSearchRepo(t *testing.T) {
	log.Println("== TestSearchRepo ==")
	c := newTestClient()

	repo, err := createTestRepo(t, "RepoSearch1", c)
	assert.NoError(t, err)
	assert.NoError(t, c.AddRepoTopic(repo.Owner.UserName, repo.Name, "TestTopic1"))
	assert.NoError(t, c.AddRepoTopic(repo.Owner.UserName, repo.Name, "TestTopic2"))

	repo, err = createTestRepo(t, "RepoSearch2", c)
	assert.NoError(t, err)
	assert.NoError(t, c.AddRepoTopic(repo.Owner.UserName, repo.Name, "TestTopic1"))

	repos, err := c.SearchRepos(SearchRepoOptions{
		Keyword:     "Search1",
		IncludeDesc: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 1)

	repos, err = c.SearchRepos(SearchRepoOptions{
		Keyword:     "Search",
		IncludeDesc: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 2)

	repos, err = c.SearchRepos(SearchRepoOptions{
		Keyword: "TestTopic1",
		Topic:   true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 2)

	repos, err = c.SearchRepos(SearchRepoOptions{
		Keyword: "TestTopic2",
		Topic:   true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 1)

	err = c.DeleteRepo(repo.Owner.UserName, repo.Name)
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
