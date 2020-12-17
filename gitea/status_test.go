// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommitStatus(t *testing.T) {
	log.Println("== TestCommitStatus ==")
	c := newTestClient()
	user, _, err := c.GetMyUserInfo()
	assert.NoError(t, err)

	var repoName = "CommitStatuses"
	origRepo, err := createTestRepo(t, repoName, c)
	if !assert.NoError(t, err) {
		return
	}

	commits, _, _ := c.ListRepoCommits(user.UserName, repoName, ListCommitOptions{
		ListOptions: ListOptions{},
		SHA:         origRepo.DefaultBranch,
	})
	if !assert.Len(t, commits, 1) {
		return
	}
	sha := commits[0].SHA

	combiStats, resp, err := c.GetCombinedStatus(user.UserName, repoName, sha)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, combiStats)
	assert.EqualValues(t, 0, combiStats.TotalCount)

	statuses, resp, err := c.ListStatuses(user.UserName, repoName, sha, ListStatusesOption{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, statuses)
	assert.Len(t, statuses, 0)

	createStatus(t, c, user.UserName, repoName, sha, "http://dummy.test", "start testing", "ultraCI", StatusPending)
	createStatus(t, c, user.UserName, repoName, sha, "https://more.secure", "just a warning", "warn/bot", StatusWarning)
	createStatus(t, c, user.UserName, repoName, sha, "http://dummy.test", "test failed", "ultraCI", StatusFailure)
	createStatus(t, c, user.UserName, repoName, sha, "http://dummy.test", "start testing", "ultraCI", StatusPending)
	createStatus(t, c, user.UserName, repoName, sha, "http://dummy.test", "test passed", "ultraCI", StatusSuccess)

	statuses, resp, err = c.ListStatuses(user.UserName, repoName, sha, ListStatusesOption{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, statuses)
	assert.Len(t, statuses, 5)

	combiStats, resp, err = c.GetCombinedStatus(user.UserName, repoName, sha)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, combiStats)
	assert.EqualValues(t, 2, combiStats.TotalCount)
	assert.EqualValues(t, StatusState("warning"), combiStats.State)
	assert.Len(t, combiStats.Statuses, 2)
}

func createStatus(t *testing.T, c *Client, userName, repoName, sha, url, desc, context string, state StatusState) {
	stats, resp, err := c.CreateStatus(userName, repoName, sha, CreateStatusOption{
		State:       state,
		TargetURL:   url,
		Description: desc,
		Context:     context,
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, stats)
	assert.EqualValues(t, state, stats.State)
}
