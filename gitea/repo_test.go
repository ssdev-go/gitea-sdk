// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRepo(t *testing.T) {
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
