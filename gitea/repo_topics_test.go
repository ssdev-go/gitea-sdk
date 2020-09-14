// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoTopics(t *testing.T) {
	log.Println("== TestRepoTopics ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "RandomTopic", c)
	assert.NoError(t, err)

	// Add
	_, err = c.AddRepoTopic(repo.Owner.UserName, repo.Name, "best")
	assert.NoError(t, err)
	_, err = c.AddRepoTopic(repo.Owner.UserName, repo.Name, "git")
	assert.NoError(t, err)
	_, err = c.AddRepoTopic(repo.Owner.UserName, repo.Name, "gitea")
	assert.NoError(t, err)
	_, err = c.AddRepoTopic(repo.Owner.UserName, repo.Name, "drone")
	assert.NoError(t, err)

	// Get List
	tl, _, err := c.ListRepoTopics(repo.Owner.UserName, repo.Name, ListRepoTopicsOptions{})
	assert.NoError(t, err)
	assert.Len(t, tl, 4)

	// Del
	_, err = c.DeleteRepoTopic(repo.Owner.UserName, repo.Name, "drone")
	assert.NoError(t, err)
	_, err = c.DeleteRepoTopic(repo.Owner.UserName, repo.Name, "best")
	assert.NoError(t, err)
	tl, _, err = c.ListRepoTopics(repo.Owner.UserName, repo.Name, ListRepoTopicsOptions{})
	assert.NoError(t, err)
	assert.Len(t, tl, 2)

	// Set List
	newTopics := []string{"analog", "digital", "cat"}
	_, err = c.SetRepoTopics(repo.Owner.UserName, repo.Name, newTopics)
	assert.NoError(t, err)
	tl, _, _ = c.ListRepoTopics(repo.Owner.UserName, repo.Name, ListRepoTopicsOptions{})
	assert.Len(t, tl, 3)

	sort.Strings(tl)
	sort.Strings(newTopics)
	assert.EqualValues(t, newTopics, tl)
}
