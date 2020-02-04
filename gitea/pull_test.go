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
		State: StateAll,
		Sort:  "leastupdate",
	})
	assert.NoError(t, err)
	assert.Len(t, pulls, 0)

	//ToDo add git stuff to have different branches witch can be used to create PRs and test merge etc ...

	// GetPullRequest get information of one PR
	//func (c *Client) GetPullRequest(owner, repo string, index int64) (*PullRequest, error)

	// CreatePullRequest create pull request with options
	//func (c *Client) CreatePullRequest(owner, repo string, opt CreatePullRequestOption) (*PullRequest, error)

	// EditPullRequest modify pull request with PR id and options
	//func (c *Client) EditPullRequest(owner, repo string, index int64, opt EditPullRequestOption) (*PullRequest, error)

	// MergePullRequest merge a PR to repository by PR id
	//func (c *Client) MergePullRequest(owner, repo string, index int64, opt MergePullRequestOption) (*MergePullRequestResponse, error)

	// IsPullRequestMerged test if one PR is merged to one repository
	//func (c *Client) IsPullRequestMerged(owner, repo string, index int64) (bool, error)
}
