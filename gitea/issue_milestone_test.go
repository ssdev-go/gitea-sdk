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

func TestMilestones(t *testing.T) {
	log.Println("== TestMilestones ==")
	c := newTestClient()

	repo, _ := createTestRepo(t, "TestMilestones", c)
	now := time.Now()
	future := time.Unix(1896134400, 0) //2030-02-01
	closed := "closed"

	// CreateMilestone 4x
	m1, err := c.CreateMilestone(repo.Owner.UserName, repo.Name, CreateMilestoneOption{Title: "v1.0", Description: "First Version", Deadline: &now})
	assert.NoError(t, err)
	_, err = c.CreateMilestone(repo.Owner.UserName, repo.Name, CreateMilestoneOption{Title: "v2.0", Description: "Second Version", Deadline: &future})
	assert.NoError(t, err)
	_, err = c.CreateMilestone(repo.Owner.UserName, repo.Name, CreateMilestoneOption{Title: "v3.0", Description: "Third Version", Deadline: nil})
	assert.NoError(t, err)
	m4, err := c.CreateMilestone(repo.Owner.UserName, repo.Name, CreateMilestoneOption{Title: "temp", Description: "part time milestone"})
	assert.NoError(t, err)

	// EditMilestone
	m1, err = c.EditMilestone(repo.Owner.UserName, repo.Name, m1.ID, EditMilestoneOption{Description: &closed, State: &closed})
	assert.NoError(t, err)

	// DeleteMilestone
	assert.NoError(t, c.DeleteMilestone(repo.Owner.UserName, repo.Name, m4.ID))

	// ListRepoMilestones
	ml, err := c.ListRepoMilestones(repo.Owner.UserName, repo.Name, ListMilestoneOption{})
	assert.NoError(t, err)
	assert.Len(t, ml, 2)
	ml, err = c.ListRepoMilestones(repo.Owner.UserName, repo.Name, ListMilestoneOption{State: "closed"})
	assert.NoError(t, err)
	assert.Len(t, ml, 1)
	ml, err = c.ListRepoMilestones(repo.Owner.UserName, repo.Name, ListMilestoneOption{State: "all"})
	assert.NoError(t, err)
	assert.Len(t, ml, 3)

	// GetMilestone
	_, err = c.GetMilestone(repo.Owner.UserName, repo.Name, m4.ID)
	assert.Error(t, err)
	m, err := c.GetMilestone(repo.Owner.UserName, repo.Name, m1.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, m1, m)
}
