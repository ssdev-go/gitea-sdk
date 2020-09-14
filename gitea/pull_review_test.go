// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPullReview(t *testing.T) {
	log.Println("== TestPullReview ==")
	c := newTestClient()

	var repoName = "Reviews"
	repo, pull, submitter, reviewer, success := preparePullReviewTest(t, c, repoName)
	if !success {
		return
	}
	defer c.AdminDeleteUser(reviewer.UserName)
	defer c.AdminDeleteUser(submitter.UserName)

	// CreatePullReview
	r1, _, err := c.CreatePullReview(repo.Owner.UserName, repo.Name, pull.Index, CreatePullReviewOptions{
		State: ReviewStateComment,
		Body:  "I'll have a look at it later",
	})
	assert.NoError(t, err)
	if assert.NotNil(t, r1) {
		assert.EqualValues(t, ReviewStateComment, r1.State)
		assert.EqualValues(t, 1, r1.Reviewer.ID)
	}

	c.SetSudo(submitter.UserName)
	r2, _, err := c.CreatePullReview(repo.Owner.UserName, repo.Name, pull.Index, CreatePullReviewOptions{
		State: ReviewStateApproved,
		Body:  "lgtm it myself",
	})
	assert.Error(t, err)
	r2, _, err = c.CreatePullReview(repo.Owner.UserName, repo.Name, pull.Index, CreatePullReviewOptions{
		State: ReviewStateComment,
		Body:  "no seriously please have a look at it",
	})
	assert.NoError(t, err)
	assert.NotNil(t, r2)

	c.SetSudo(reviewer.UserName)
	r3, _, err := c.CreatePullReview(repo.Owner.UserName, repo.Name, pull.Index, CreatePullReviewOptions{
		State: ReviewStateApproved,
		Body:  "lgtm",
		Comments: []CreatePullReviewComment{{
			Path:       "WOW-file",
			Body:       "no better name - really?",
			NewLineNum: 1,
		},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, r3)

	// ListPullReviews
	c.SetSudo("")
	rl, _, err := c.ListPullReviews(repo.Owner.UserName, repo.Name, pull.Index, ListPullReviewsOptions{})
	if !assert.NoError(t, err) {
		return
	}
	assert.Len(t, rl, 3)
	for i := range rl {
		assert.EqualValues(t, pull.HTMLURL, rl[i].HTMLPullURL)
		if rl[i].CodeCommentsCount == 1 {
			assert.EqualValues(t, reviewer.ID, rl[i].Reviewer.ID)
		}
	}

	// GetPullReview
	rNew, _, err := c.GetPullReview(repo.Owner.UserName, repo.Name, pull.Index, r3.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, r3, rNew)

	// DeletePullReview
	c.SetSudo(submitter.UserName)
	_, err = c.DeletePullReview(repo.Owner.UserName, repo.Name, pull.Index, r2.ID)
	assert.NoError(t, err)
	_, err = c.DeletePullReview(repo.Owner.UserName, repo.Name, pull.Index, r3.ID)
	assert.Error(t, err)

	// SubmitPullReview
	c.SetSudo("")
	r4, _, err := c.CreatePullReview(repo.Owner.UserName, repo.Name, pull.Index, CreatePullReviewOptions{
		Body: "...",
		Comments: []CreatePullReviewComment{{
			Path:       "WOW-file",
			Body:       "its ok",
			NewLineNum: 1,
		},
		},
	})
	assert.NoError(t, err)
	r5, _, err := c.CreatePullReview(repo.Owner.UserName, repo.Name, pull.Index, CreatePullReviewOptions{
		Body: "...",
		Comments: []CreatePullReviewComment{{
			Path:       "WOW-file",
			Body:       "hehe and here it is",
			NewLineNum: 3,
		},
		},
	})
	assert.NoError(t, err)
	assert.EqualValues(t, r4.ID, r5.ID)

	r, _, err := c.SubmitPullReview(repo.Owner.UserName, repo.Name, pull.Index, r4.ID, SubmitPullReviewOptions{
		State: ReviewStateRequestChanges,
		Body:  "one nit",
	})
	assert.NoError(t, err)
	assert.EqualValues(t, r4.ID, r.ID)
	assert.EqualValues(t, ReviewStateRequestChanges, r.State)

	// ListPullReviewComments
	rcl, _, err := c.ListPullReviewComments(repo.Owner.UserName, repo.Name, pull.Index, r.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, r.CodeCommentsCount, len(rcl))
	for _, rc := range rcl {
		assert.EqualValues(t, pull.HTMLURL, rc.HTMLPullURL)
		if rc.LineNum == 3 {
			assert.EqualValues(t, "hehe and here it is", rc.Body)
		} else {
			assert.EqualValues(t, 1, rc.LineNum)
			assert.EqualValues(t, "its ok", rc.Body)
		}
	}

}

func preparePullReviewTest(t *testing.T, c *Client, repoName string) (*Repository, *PullRequest, *User, *User, bool) {
	repo, err := createTestRepo(t, repoName, c)
	if !assert.NoError(t, err) {
		return nil, nil, nil, nil, false
	}

	pullSubmitter := createTestUser(t, "pull_submitter", c)
	write := AccessModeWrite
	_, err = c.AddCollaborator(repo.Owner.UserName, repo.Name, pullSubmitter.UserName, AddCollaboratorOption{
		Permission: &write,
	})
	assert.NoError(t, err)

	c.SetSudo("pull_submitter")

	newFile, _, err := c.CreateFile(repo.Owner.UserName, repo.Name, "WOW-file", CreateFileOptions{
		Content: "QSBuZXcgRmlsZQoKYW5kIHNvbWUgbGluZXMK",
		FileOptions: FileOptions{
			Message:       "creat a new file",
			BranchName:    "master",
			NewBranchName: "new_file",
		},
	})

	if !assert.NoError(t, err) || !assert.NotNil(t, newFile) {
		return nil, nil, nil, nil, false
	}

	pull, _, err := c.CreatePullRequest(c.username, repoName, CreatePullRequestOption{
		Base:  "master",
		Head:  "new_file",
		Title: "Creat a NewFile",
	})
	assert.NoError(t, err)
	assert.NotNil(t, pull)

	c.SetSudo("")

	reviewer := createTestUser(t, "pull_reviewer", c)
	admin := AccessModeAdmin
	_, err = c.AddCollaborator(repo.Owner.UserName, repo.Name, pullSubmitter.UserName, AddCollaboratorOption{
		Permission: &admin,
	})
	assert.NoError(t, err)

	return repo, pull, pullSubmitter, reviewer, pull.Poster.ID == pullSubmitter.ID
}
