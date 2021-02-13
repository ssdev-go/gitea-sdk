// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	log.Println("== TestTags ==")
	c := newTestClient()

	repo, _ := createTestRepo(t, "TestTags", c)

	// Create Tags
	createTestTag(t, c, repo, "tag1")

	tags, _, err := c.ListRepoTags(repo.Owner.UserName, repo.Name, ListRepoTagsOptions{})
	assert.NoError(t, err)
	assert.Len(t, tags, 1)

	// DeleteReleaseTag
	resp, err := c.DeleteTag(repo.Owner.UserName, repo.Name, "tag1")
	assert.NoError(t, err)
	assert.EqualValues(t, 204, resp.StatusCode)
	tags, _, err = c.ListRepoTags(repo.Owner.UserName, repo.Name, ListRepoTagsOptions{})
	assert.NoError(t, err)
	assert.Len(t, tags, 0)
}

// createTestTag use create release api since there exist no api to create tag only
// https://github.com/go-gitea/gitea/issues/14669
func createTestTag(t *testing.T, c *Client, repo *Repository, name string) {
	rel, _, err := c.CreateRelease(repo.Owner.UserName, repo.Name, CreateReleaseOption{
		TagName: name,
		Target:  "master",
		Title:   "TMP Release",
	})
	assert.NoError(t, err)
	_, err = c.DeleteRelease(repo.Owner.UserName, repo.Name, rel.ID)
	assert.NoError(t, err)
}
