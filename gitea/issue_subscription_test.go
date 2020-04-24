// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIssue is main func witch call all Tests for Issue API
// (to make sure they are on correct order)
func TestIssueSubscription(t *testing.T) {
	log.Println("== TestCreateIssues ==")

	c := newTestClient()
	repo, _ := createTestRepo(t, "IssueWatch", c)
	createTestIssue(t, c, repo.Name, "First Issue", "", nil, nil, 0, nil, false, false)

	wi, err := c.CheckIssueSubscription(repo.Owner.UserName, repo.Name, 1)
	assert.NoError(t, err)
	assert.True(t, wi.Subscribed)

	assert.NoError(t, c.UnWatchRepo(repo.Owner.UserName, repo.Name))
	wi, err = c.CheckIssueSubscription(repo.Owner.UserName, repo.Name, 1)
	assert.NoError(t, err)
	assert.True(t, wi.Subscribed)

	assert.NoError(t, c.IssueSubscribe(repo.Owner.UserName, repo.Name, 1))
	wi, err = c.CheckIssueSubscription(repo.Owner.UserName, repo.Name, 1)
	assert.NoError(t, err)
	assert.True(t, wi.Subscribed)

	assert.NoError(t, c.IssueUnSubscribe(repo.Owner.UserName, repo.Name, 1))
	wi, err = c.CheckIssueSubscription(repo.Owner.UserName, repo.Name, 1)
	assert.NoError(t, err)
	assert.False(t, wi.Subscribed)

	assert.NoError(t, c.WatchRepo(repo.Owner.UserName, repo.Name))
	wi, err = c.CheckIssueSubscription(repo.Owner.UserName, repo.Name, 1)
	assert.NoError(t, err)
	assert.False(t, wi.Subscribed)
}
