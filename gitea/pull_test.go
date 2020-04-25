// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPull(t *testing.T) {
	log.Println("== TestPull ==")
	c := newTestClient()
	user, err := c.GetMyUserInfo()
	assert.NoError(t, err)

	var repoName = "repo_pull_test"
	_, err = createTestRepo(t, repoName, c)
	if err != nil {
		return
	}

	// ListRepoPullRequests list PRs of one repository
	pulls, err := c.ListRepoPullRequests(user.UserName, repoName, ListPullRequestsOptions{
		Page:  1,
		State: "all",
		Sort:  "leastupdate",
	})
	assert.NoError(t, err)
	assert.Len(t, pulls, 0)
}
