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

// TestIssue is main func witch call all Tests for Issue API
// (to make sure they are on correct order)
func TestIssue(t *testing.T) {
	c := newTestClient()

	createIssue(t, c)
	listIssues(t, c)
}

func createIssue(t *testing.T, c *Client) {
	log.Println("== TestCreateIssues ==")

	user, err := c.GetMyUserInfo()
	assert.NoError(t, err)
	repo, _ := createTestRepo(t, "IssueTestsRepo", c)

	createOne := func(title, body string, assignees []string, deadline *time.Time, milestone int64, labels []int64, closed, shouldFail bool) {
		issue, e := c.CreateIssue(user.UserName, repo.Name, CreateIssueOption{
			Title:     title,
			Body:      body,
			Assignees: assignees,
			Deadline:  deadline,
			Milestone: milestone,
			Labels:    labels,
			Closed:    closed,
		})
		if shouldFail {
			assert.Error(t, e)
			return
		}
		assert.NoError(t, e)
		assert.NotEmpty(t, issue)
		assert.EqualValues(t, title, issue.Title)
		assert.EqualValues(t, body, issue.Body)
		assert.EqualValues(t, len(assignees), len(issue.Assignees))
		for i, a := range issue.Assignees {
			assert.EqualValues(t, assignees[i], a.UserName)
		}
		if milestone > 0 {
			assert.EqualValues(t, milestone, issue.Milestone.ID)
		}
		assert.EqualValues(t, len(labels), len(issue.Labels))
		if closed {
			assert.False(t, issue.Closed.IsZero())
		} else {
			assert.Empty(t, issue.Closed)
		}
	}

	nowTime := time.Now()
	mile, _ := c.CreateMilestone(user.UserName, repo.Name, CreateMilestoneOption{Title: "mile1"})
	label1, _ := c.CreateLabel(user.UserName, repo.Name, CreateLabelOption{Name: "Label1", Description: "a", Color: "#ee0701"})
	label2, _ := c.CreateLabel(user.UserName, repo.Name, CreateLabelOption{Name: "Label2", Description: "b", Color: "#128a0c"})

	createOne("First Issue", "", nil, nil, 0, nil, false, false)
	createOne("Issue 2", "closed isn't it?", nil, nil, 0, nil, true, false)
	createOne("Issue 3", "", nil, nil, 0, nil, true, false)
	createOne("Feature: spam protect 4", "explain explain explain", []string{user.UserName}, &nowTime, 0, nil, true, false)
	createOne("W 123", "", nil, &nowTime, mile.ID, nil, false, false)
	createOne("First Issue", "", nil, nil, 0, nil, false, false)
	createOne("Do it soon!", "is important!", []string{user.UserName}, &nowTime, mile.ID, []int64{label1.ID, label2.ID}, false, false)
	createOne("Job Done", "you never know", nil, nil, mile.ID, []int64{label2.ID}, true, false)
	createOne("", "you never know", nil, nil, mile.ID, nil, true, true)
}

func listIssues(t *testing.T, c *Client) {
	log.Println("== TestListIssues ==")

	issues, err := c.ListRepoIssues("test01", "IssueTestsRepo", ListIssueOption{
		Labels:  []string{"Label2"},
		KeyWord: "",
		State:   "all",
	})
	assert.NoError(t, err)
	assert.Len(t, issues, 2)

	issues, err = c.ListIssues(ListIssueOption{
		Labels:  []string{"Label2"},
		KeyWord: "Done",
		State:   "",
	})
	assert.NoError(t, err)
	assert.Len(t, issues, 1)

	issues, err = c.ListRepoIssues("test01", "IssueTestsRepo", ListIssueOption{})
	assert.NoError(t, err)
	assert.Len(t, issues, 4)
}
