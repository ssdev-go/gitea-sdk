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
	var forkOrg = "ForkOrg"
	if !preparePullTest(t, c, repoName, forkOrg) {
		return
	}

	// ListRepoPullRequests list PRs of one repository
	pulls, err := c.ListRepoPullRequests(user.UserName, repoName, ListPullRequestsOptions{})
	assert.NoError(t, err)
	assert.Len(t, pulls, 0)

	pullUpdateFile, err := c.CreatePullRequest(c.username, repoName, CreatePullRequestOption{
		Base:  "master",
		Head:  forkOrg + ":overwrite_licence",
		Title: "overwrite a file",
	})
	assert.NoError(t, err)
	assert.NotNil(t, pullUpdateFile)

	pullNewFile, err := c.CreatePullRequest(c.username, repoName, CreatePullRequestOption{
		Base:  "master",
		Head:  forkOrg + ":new_file",
		Title: "create a file",
	})
	assert.NoError(t, err)
	assert.NotNil(t, pullNewFile)

	pullConflict, err := c.CreatePullRequest(c.username, repoName, CreatePullRequestOption{
		Base:  "master",
		Head:  forkOrg + ":will_conflict",
		Title: "this pull will conflict",
	})
	assert.NoError(t, err)
	assert.NotNil(t, pullConflict)

	pulls, err = c.ListRepoPullRequests(user.UserName, repoName, ListPullRequestsOptions{})
	assert.NoError(t, err)
	assert.Len(t, pulls, 3)

	diff, err := c.GetPullRequestDiff(c.username, repoName, pullUpdateFile.Index)
	assert.NoError(t, err)
	assert.Len(t, diff, 1310)
	patch, err := c.GetPullRequestPatch(c.username, repoName, pullUpdateFile.Index)
	assert.NoError(t, err)
	assert.True(t, len(patch) > len(diff))

	// test Update pull
	pr, err := c.GetPullRequest(user.UserName, repoName, pullUpdateFile.Index)
	assert.NoError(t, err)
	assert.False(t, pullUpdateFile.HasMerged)
	assert.True(t, pullUpdateFile.Mergeable)
	merged, err := c.MergePullRequest(user.UserName, repoName, pullUpdateFile.Index, MergePullRequestOption{
		Style:   MergeStyleSquash,
		Title:   pullUpdateFile.Title,
		Message: "squash: " + pullUpdateFile.Title,
	})
	assert.NoError(t, err)
	assert.True(t, merged)
	merged, err = c.IsPullRequestMerged(user.UserName, repoName, pullUpdateFile.Index)
	assert.NoError(t, err)
	assert.True(t, merged)
	pr, err = c.GetPullRequest(user.UserName, repoName, pullUpdateFile.Index)
	assert.NoError(t, err)
	assert.EqualValues(t, pullUpdateFile.Head.Name, pr.Head.Name)
	assert.EqualValues(t, pullUpdateFile.Base.Name, pr.Base.Name)
	assert.NotEqual(t, pullUpdateFile.Base.Sha, pr.Base.Sha)
	assert.Len(t, *pr.MergedCommitID, 40)
	assert.True(t, pr.HasMerged)

	// test conflict pull
	pr, err = c.GetPullRequest(user.UserName, repoName, pullConflict.Index)
	assert.NoError(t, err)
	assert.False(t, pullConflict.HasMerged)
	assert.False(t, pullConflict.Mergeable)
	merged, err = c.MergePullRequest(user.UserName, repoName, pullConflict.Index, MergePullRequestOption{
		Style:   MergeStyleMerge,
		Title:   "pullConflict",
		Message: "pullConflict Msg",
	})
	assert.NoError(t, err)
	assert.False(t, merged)
	merged, err = c.IsPullRequestMerged(user.UserName, repoName, pullConflict.Index)
	assert.NoError(t, err)
	assert.False(t, merged)
	pr, err = c.GetPullRequest(user.UserName, repoName, pullConflict.Index)
	assert.NoError(t, err)
	assert.Nil(t, pr.MergedCommitID)
	assert.False(t, pr.HasMerged)

	state := StateClosed
	pr, err = c.EditPullRequest(user.UserName, repoName, pullConflict.Index, EditPullRequestOption{
		Title: "confl",
		State: &state,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, state, pr.State)

	pulls, err = c.ListRepoPullRequests(user.UserName, repoName, ListPullRequestsOptions{
		State: StateClosed,
		Sort:  "leastupdate",
	})
	assert.NoError(t, err)
	assert.Len(t, pulls, 2)
}

func preparePullTest(t *testing.T, c *Client, repoName, forkOrg string) bool {
	_ = c.DeleteRepo(forkOrg, repoName)
	_ = c.DeleteRepo(c.username, repoName)
	_ = c.DeleteOrg(forkOrg)

	origRepo, err := createTestRepo(t, repoName, c)
	if !assert.NoError(t, err) {
		return false
	}
	org, err := c.CreateOrg(CreateOrgOption{Name: forkOrg})
	assert.NoError(t, err)
	forkRepo, err := c.CreateFork(origRepo.Owner.UserName, origRepo.Name, CreateForkOption{Organization: &org.UserName})
	assert.NoError(t, err)
	assert.NotNil(t, forkRepo)

	masterLicence, err := c.GetContents(forkRepo.Owner.UserName, forkRepo.Name, "master", "LICENSE")
	if !assert.NoError(t, err) || !assert.NotNil(t, masterLicence) {
		return false
	}

	updatedFile, err := c.UpdateFile(forkRepo.Owner.UserName, forkRepo.Name, "LICENSE", UpdateFileOptions{
		FileOptions: FileOptions{
			Message:       "Overwrite",
			BranchName:    "master",
			NewBranchName: "overwrite_licence",
		},
		SHA:     masterLicence.SHA,
		Content: "Tk9USElORyBJUyBIRVJFIEFOWU1PUkUKSUYgWU9VIExJS0UgVE8gRklORCBTT01FVEhJTkcKV0FJVCBGT1IgVEhFIEZVVFVSRQo=",
	})
	if !assert.NoError(t, err) || !assert.NotNil(t, updatedFile) {
		return false
	}

	newFile, err := c.CreateFile(forkRepo.Owner.UserName, forkRepo.Name, "WOW-file", CreateFileOptions{
		Content: "QSBuZXcgRmlsZQo=",
		FileOptions: FileOptions{
			Message:       "creat a new file",
			BranchName:    "master",
			NewBranchName: "new_file",
		},
	})
	if !assert.NoError(t, err) || !assert.NotNil(t, newFile) {
		return false
	}

	conflictFile1, err := c.CreateFile(origRepo.Owner.UserName, origRepo.Name, "bad-file", CreateFileOptions{
		Content: "U3RhcnQgQ29uZmxpY3QK",
		FileOptions: FileOptions{
			Message:    "Start Conflict",
			BranchName: "master",
		},
	})
	if !assert.NoError(t, err) || !assert.NotNil(t, conflictFile1) {
		return false
	}

	conflictFile2, err := c.CreateFile(forkRepo.Owner.UserName, forkRepo.Name, "bad-file", CreateFileOptions{
		Content: "V2lsbEhhdmUgQ29uZmxpY3QK",
		FileOptions: FileOptions{
			Message:       "creat a new file witch will conflict",
			BranchName:    "master",
			NewBranchName: "will_conflict",
		},
	})
	if !assert.NoError(t, err) || !assert.NotNil(t, conflictFile2) {
		return false
	}

	return true
}
