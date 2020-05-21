// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoBranches(t *testing.T) {
	log.Println("== TestRepoBranches ==")
	c := newTestClient()
	var repoName = "branches"

	repo := prepareBranchTest(t, c, repoName)
	if repo == nil {
		return
	}

	bl, err := c.ListRepoBranches(repo.Owner.UserName, repo.Name, ListRepoBranchesOptions{})
	assert.NoError(t, err)
	assert.Len(t, bl, 3)
	assert.EqualValues(t, "feature", bl[0].Name)
	assert.EqualValues(t, "master", bl[1].Name)
	assert.EqualValues(t, "update", bl[2].Name)

	b, err := c.GetRepoBranch(repo.Owner.UserName, repo.Name, "update")
	assert.NoError(t, err)
	assert.EqualValues(t, bl[2].Commit.ID, b.Commit.ID)
	assert.EqualValues(t, bl[2].Commit.Added, b.Commit.Added)

	s, err := c.DeleteRepoBranch(repo.Owner.UserName, repo.Name, "master")
	assert.NoError(t, err)
	assert.False(t, s)
	s, err = c.DeleteRepoBranch(repo.Owner.UserName, repo.Name, "feature")
	assert.NoError(t, err)
	assert.True(t, s)

	bl, err = c.ListRepoBranches(repo.Owner.UserName, repo.Name, ListRepoBranchesOptions{})
	assert.NoError(t, err)
	assert.Len(t, bl, 2)

	b, err = c.GetRepoBranch(repo.Owner.UserName, repo.Name, "feature")
	assert.Error(t, err)
	assert.Nil(t, b)
}

func TestRepoBranchProtection(t *testing.T) {
	log.Println("== TestRepoBranchProtection ==")
	c := newTestClient()
	var repoName = "BranchProtection"

	repo := prepareBranchTest(t, c, repoName)
	if repo == nil {
		return
	}
	assert.NotNil(t, repo)

	// ListBranchProtections
	bpl, err := c.ListBranchProtections(repo.Owner.UserName, repo.Name, ListBranchProtectionsOptions{})
	assert.NoError(t, err)
	assert.Len(t, bpl, 0)

	// CreateBranchProtection
	bp, err := c.CreateBranchProtection(repo.Owner.UserName, repo.Name, CreateBranchProtectionOption{
		BranchName:              "master",
		EnablePush:              true,
		EnablePushWhitelist:     true,
		PushWhitelistUsernames:  []string{"test01"},
		EnableMergeWhitelist:    true,
		MergeWhitelistUsernames: []string{"test01"},
		BlockOnOutdatedBranch:   true,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, "master", bp.BranchName)
	assert.EqualValues(t, false, bp.EnableStatusCheck)
	assert.EqualValues(t, true, bp.EnablePush)
	assert.EqualValues(t, true, bp.EnablePushWhitelist)
	assert.EqualValues(t, []string{"test01"}, bp.PushWhitelistUsernames)

	bp, err = c.CreateBranchProtection(repo.Owner.UserName, repo.Name, CreateBranchProtectionOption{
		BranchName:              "update",
		EnablePush:              false,
		EnableMergeWhitelist:    true,
		MergeWhitelistUsernames: []string{"test01"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, bp)

	bpl, err = c.ListBranchProtections(repo.Owner.UserName, repo.Name, ListBranchProtectionsOptions{})
	assert.NoError(t, err)
	assert.Len(t, bpl, 2)

	// GetBranchProtection
	bp, err = c.GetBranchProtection(repo.Owner.UserName, repo.Name, bpl[0].BranchName)
	assert.NoError(t, err)
	assert.EqualValues(t, bpl[0], bp)

	optTrue := true
	optFalse := false
	one := int64(1)

	// EditBranchProtection
	bp, err = c.EditBranchProtection(repo.Owner.UserName, repo.Name, bpl[0].BranchName, EditBranchProtectionOption{
		EnablePush:                  &optFalse,
		EnablePushWhitelist:         &optFalse,
		PushWhitelistUsernames:      nil,
		RequiredApprovals:           &one,
		EnableApprovalsWhitelist:    &optTrue,
		ApprovalsWhitelistUsernames: []string{"test01"},
	})
	assert.NoError(t, err)
	assert.NotEqual(t, bpl[0], bp)
	assert.EqualValues(t, bpl[0].BranchName, bp.BranchName)
	assert.EqualValues(t, bpl[0].EnableMergeWhitelist, bp.EnableMergeWhitelist)
	assert.EqualValues(t, bpl[0].Created, bp.Created)

	// DeleteBranchProtection
	err = c.DeleteBranchProtection(repo.Owner.UserName, repo.Name, bpl[1].BranchName)
	assert.NoError(t, err)
	bpl, err = c.ListBranchProtections(repo.Owner.UserName, repo.Name, ListBranchProtectionsOptions{})
	assert.NoError(t, err)
	assert.Len(t, bpl, 1)
}

func prepareBranchTest(t *testing.T, c *Client, repoName string) *Repository {
	origRepo, err := createTestRepo(t, repoName, c)
	if !assert.NoError(t, err) {
		return nil
	}

	masterLicence, err := c.GetContents(origRepo.Owner.UserName, origRepo.Name, "master", "README.md")
	if !assert.NoError(t, err) || !assert.NotNil(t, masterLicence) {
		return nil
	}

	updatedFile, err := c.UpdateFile(origRepo.Owner.UserName, origRepo.Name, "README.md", UpdateFileOptions{
		FileOptions: FileOptions{
			Message:       "update it",
			BranchName:    "master",
			NewBranchName: "update",
		},
		SHA:     masterLicence.SHA,
		Content: "Tk9USElORyBJUyBIRVJFIEFOWU1PUkUKSUYgWU9VIExJS0UgVE8gRklORCBTT01FVEhJTkcKV0FJVCBGT1IgVEhFIEZVVFVSRQo=",
	})
	if !assert.NoError(t, err) || !assert.NotNil(t, updatedFile) {
		return nil
	}

	newFile, err := c.CreateFile(origRepo.Owner.UserName, origRepo.Name, "WOW-file", CreateFileOptions{
		Content: "QSBuZXcgRmlsZQo=",
		FileOptions: FileOptions{
			Message:       "creat a new file",
			BranchName:    "master",
			NewBranchName: "feature",
		},
	})
	if !assert.NoError(t, err) || !assert.NotNil(t, newFile) {
		return nil
	}

	return origRepo
}
