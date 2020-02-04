// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoTransfer(t *testing.T) {
	log.Printf("== TestRepoTransfer ==")
	c := newTestClient()

	org, err := c.AdminCreateOrg(c.username, CreateOrgOption{UserName: "TransferOrg"})
	assert.NoError(t, err)
	repo, err := createTestRepo(t, "ToMove", c)
	assert.NoError(t, err)

	newRepo, err := c.TransferRepo(repo.Owner.UserName, repo.Name, TransferRepoOption{NewOwner: org.UserName})
	assert.NoError(t, err)
	assert.NotNil(t, newRepo)

	repo, err = createTestRepo(t, "ToMove", c)
	assert.NoError(t, err)
	_, err = c.TransferRepo(repo.Owner.UserName, repo.Name, TransferRepoOption{NewOwner: org.UserName})
	assert.Error(t, err)

	assert.NoError(t, c.DeleteRepo(repo.Owner.UserName, repo.Name))
	assert.NoError(t, c.DeleteRepo(newRepo.Owner.UserName, newRepo.Name))
	assert.NoError(t, c.DeleteOrg(org.UserName))
}
