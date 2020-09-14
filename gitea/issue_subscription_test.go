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
	log.Println("== TestIssueSubscription ==")

	c := newTestClient()
	repo, _ := createTestRepo(t, "IssueWatch", c)
	createTestIssue(t, c, repo.Name, "First Issue", "", nil, nil, 0, nil, false, false)

	wi, _, err := c.CheckIssueSubscription(repo.Owner.UserName, repo.Name, 1)
	assert.NoError(t, err)
	assert.True(t, wi.Subscribed)

	_, err = c.UnWatchRepo(repo.Owner.UserName, repo.Name)
	assert.NoError(t, err)
	wi, _, err = c.CheckIssueSubscription(repo.Owner.UserName, repo.Name, 1)
	assert.NoError(t, err)
	assert.True(t, wi.Subscribed)

	_, err = c.IssueSubscribe(repo.Owner.UserName, repo.Name, 1)
	if assert.Error(t, err) {
		assert.EqualValues(t, "already subscribed", err.Error())
	}
	wi, _, err = c.CheckIssueSubscription(repo.Owner.UserName, repo.Name, 1)
	assert.NoError(t, err)
	assert.True(t, wi.Subscribed)

	_, err = c.IssueUnSubscribe(repo.Owner.UserName, repo.Name, 1)
	assert.NoError(t, err)
	wi, _, err = c.CheckIssueSubscription(repo.Owner.UserName, repo.Name, 1)
	assert.NoError(t, err)
	assert.False(t, wi.Subscribed)

	_, err = c.WatchRepo(repo.Owner.UserName, repo.Name)
	assert.NoError(t, err)
	wi, _, err = c.CheckIssueSubscription(repo.Owner.UserName, repo.Name, 1)
	assert.NoError(t, err)
	assert.False(t, wi.Subscribed)
}
