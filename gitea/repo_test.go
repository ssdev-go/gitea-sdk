// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepo(t *testing.T) {
	log.Println("== TestCreateRepo ==")
	c := newTestClient()
	user, _, err := c.GetMyUserInfo()
	assert.NoError(t, err)

	var repoName = "test1"
	_, _, err = c.GetRepo(user.UserName, repoName)
	if err != nil {
		repo, _, err := c.CreateRepo(CreateRepoOption{
			Name: repoName,
		})
		assert.NoError(t, err)
		assert.NotNil(t, repo)
	}

	_, err = c.DeleteRepo(user.UserName, repoName)
	assert.NoError(t, err)
}

func TestRepoMigrateAndLanguages(t *testing.T) {
	log.Println("== TestMigrateRepo ==")
	c := newTestClient()
	user, _, uErr := c.GetMyUserInfo()
	assert.NoError(t, uErr)
	_, _, err := c.GetRepo(user.UserName, "sdk-mirror")
	if err == nil {
		_, _ = c.DeleteRepo(user.UserName, "sdk-mirror")
	}

	repoM, _, err := c.MigrateRepo(MigrateRepoOption{
		CloneAddr:   "https://gitea.com/gitea/go-sdk.git",
		RepoName:    "sdk-mirror",
		RepoOwner:   user.UserName,
		Mirror:      true,
		Private:     false,
		Description: "mirror sdk",
	})
	assert.NoError(t, err)

	repoG, _, err := c.GetRepo(repoM.Owner.UserName, repoM.Name)
	assert.NoError(t, err)
	assert.EqualValues(t, repoM.ID, repoG.ID)
	assert.EqualValues(t, "master", repoG.DefaultBranch)
	assert.True(t, repoG.Mirror)
	assert.False(t, repoG.Empty)
	assert.EqualValues(t, 1, repoG.Watchers)

	log.Println("== TestRepoLanguages ==")
	time.Sleep(time.Second)
	lang, _, err := c.GetRepoLanguages(repoM.Owner.UserName, repoM.Name)
	assert.NoError(t, err)
	assert.Len(t, lang, 2)
	assert.True(t, 217441 < lang["Go"])
	assert.True(t, 3614 < lang["Makefile"] && 6000 > lang["Makefile"])
}

func TestSearchRepo(t *testing.T) {
	log.Println("== TestSearchRepo ==")
	c := newTestClient()

	repo, err := createTestRepo(t, "RepoSearch1", c)
	assert.NoError(t, err)
	_, err = c.AddRepoTopic(repo.Owner.UserName, repo.Name, "TestTopic1")
	assert.NoError(t, err)
	_, err = c.AddRepoTopic(repo.Owner.UserName, repo.Name, "TestTopic2")
	assert.NoError(t, err)

	repo, err = createTestRepo(t, "RepoSearch2", c)
	assert.NoError(t, err)
	_, err = c.AddRepoTopic(repo.Owner.UserName, repo.Name, "TestTopic1")
	assert.NoError(t, err)

	repos, _, err := c.SearchRepos(SearchRepoOptions{
		Keyword:              "Search1",
		KeywordInDescription: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 1)

	repos, _, err = c.SearchRepos(SearchRepoOptions{
		Keyword:              "Search",
		KeywordInDescription: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 2)

	repos, _, err = c.SearchRepos(SearchRepoOptions{
		Keyword:              "TestTopic1",
		KeywordInDescription: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 2)

	repos, _, err = c.SearchRepos(SearchRepoOptions{
		Keyword:              "TestTopic2",
		KeywordInDescription: true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, repos)
	assert.Len(t, repos, 1)

	_, err = c.DeleteRepo(repo.Owner.UserName, repo.Name)
	assert.NoError(t, err)
}

func TestDeleteRepo(t *testing.T) {
	log.Println("== TestDeleteRepo ==")
	c := newTestClient()
	repo, _ := createTestRepo(t, "TestDeleteRepo", c)
	_, err := c.DeleteRepo(repo.Owner.UserName, repo.Name)
	assert.NoError(t, err)
}

func TestGetArchive(t *testing.T) {
	log.Println("== TestGetArchive ==")
	c := newTestClient()
	repo, _ := createTestRepo(t, "ToDownload", c)
	time.Sleep(time.Second / 2)
	archive, _, err := c.GetArchive(repo.Owner.UserName, repo.Name, "master", ZipArchive)
	assert.NoError(t, err)
	assert.EqualValues(t, 1602, len(archive))
}

// standard func to create a init repo for test routines
func createTestRepo(t *testing.T, name string, c *Client) (*Repository, error) {
	user, _, uErr := c.GetMyUserInfo()
	assert.NoError(t, uErr)
	_, _, err := c.GetRepo(user.UserName, name)
	if err == nil {
		_, _ = c.DeleteRepo(user.UserName, name)
	}
	repo, _, err := c.CreateRepo(CreateRepoOption{
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
