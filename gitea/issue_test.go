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
	// Little sleep in order to give some time for gitea to properly store all information on database. Without this sleep, CI is a bit unstable
	time.Sleep(100 * time.Millisecond)
	editIssues(t, c)
	listIssues(t, c)
}

func createIssue(t *testing.T, c *Client) {
	log.Println("== TestCreateIssues ==")

	user, _, err := c.GetMyUserInfo()
	assert.NoError(t, err)
	repo, _ := createTestRepo(t, "IssueTestsRepo", c)

	nowTime := time.Now()
	mile, _, _ := c.CreateMilestone(user.UserName, repo.Name, CreateMilestoneOption{Title: "mile1"})
	label1, _, _ := c.CreateLabel(user.UserName, repo.Name, CreateLabelOption{Name: "Label1", Description: "a", Color: "#ee0701"})
	label2, _, _ := c.CreateLabel(user.UserName, repo.Name, CreateLabelOption{Name: "Label2", Description: "b", Color: "#128a0c"})

	createTestIssue(t, c, repo.Name, "First Issue", "", nil, nil, 0, nil, false, false)
	createTestIssue(t, c, repo.Name, "Issue 2", "closed isn't it?", nil, nil, 0, nil, true, false)
	createTestIssue(t, c, repo.Name, "Issue 3", "", nil, nil, 0, nil, true, false)
	createTestIssue(t, c, repo.Name, "Feature: spam protect 4", "explain explain explain", []string{user.UserName}, &nowTime, 0, nil, true, false)
	createTestIssue(t, c, repo.Name, "W 123", "", nil, &nowTime, mile.ID, nil, false, false)
	createTestIssue(t, c, repo.Name, "First Issue", "", nil, nil, 0, nil, false, false)
	createTestIssue(t, c, repo.Name, "Do it soon!", "is important!", []string{user.UserName}, &nowTime, mile.ID, []int64{label1.ID, label2.ID}, false, false)
	createTestIssue(t, c, repo.Name, "Job Done", "you never know", nil, nil, mile.ID, []int64{label2.ID}, true, false)
	createTestIssue(t, c, repo.Name, "", "you never know", nil, nil, mile.ID, nil, true, true)
}

func editIssues(t *testing.T, c *Client) {
	log.Println("== TestEditIssues ==")
	il, _, err := c.ListIssues(ListIssueOption{KeyWord: "soon"})
	assert.NoError(t, err)
	issue, _, err := c.GetIssue(il[0].Poster.UserName, il[0].Repository.Name, il[0].Index)
	assert.NoError(t, err)

	state := StateClosed
	issueNew, _, err := c.EditIssue(issue.Poster.UserName, issue.Repository.Name, issue.Index, EditIssueOption{
		Title: "Edited",
		Body:  OptionalString("123 test and go"),
		State: &state,
		Ref:   OptionalString("master"),
	})
	assert.NoError(t, err)
	assert.EqualValues(t, issue.ID, issueNew.ID)
	assert.EqualValues(t, "123 test and go", issueNew.Body)
	assert.EqualValues(t, "Edited", issueNew.Title)
	assert.EqualValues(t, "master", issueNew.Ref)
}

func listIssues(t *testing.T, c *Client) {
	log.Println("== TestListIssues ==")

	issues, _, err := c.ListRepoIssues("test01", "IssueTestsRepo", ListIssueOption{
		Labels:  []string{"Label1", "Label2"},
		KeyWord: "",
		State:   "all",
	})
	assert.NoError(t, err)
	assert.Len(t, issues, 1)

	issues, _, err = c.ListIssues(ListIssueOption{
		Labels:  []string{"Label2"},
		KeyWord: "Done",
		State:   "all",
	})
	assert.NoError(t, err)
	assert.Len(t, issues, 1)

	issues, _, err = c.ListRepoIssues("test01", "IssueTestsRepo", ListIssueOption{
		Milestones: []string{"mile1"},
		State:      "all",
	})
	assert.NoError(t, err)
	assert.Len(t, issues, 3)
	for i := range issues {
		if assert.NotNil(t, issues[i].Milestone) {
			assert.EqualValues(t, "mile1", issues[i].Milestone.Title)
		}
	}

	issues, _, err = c.ListRepoIssues("test01", "IssueTestsRepo", ListIssueOption{})
	assert.NoError(t, err)
	assert.Len(t, issues, 3)
}

func createTestIssue(t *testing.T, c *Client, repoName, title, body string, assignees []string, deadline *time.Time, milestone int64, labels []int64, closed, shouldFail bool) {
	user, _, err := c.GetMyUserInfo()
	assert.NoError(t, err)
	issue, _, e := c.CreateIssue(user.UserName, repoName, CreateIssueOption{
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
