// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelease(t *testing.T) {
	log.Println("== TestRelease ==")
	c := newTestClient()

	repo, _ := createTestRepo(t, "ReleaseTests", c)

	// ListReleases
	rl, _, err := c.ListReleases(repo.Owner.UserName, repo.Name, ListReleasesOptions{})
	assert.NoError(t, err)
	assert.Len(t, rl, 0)

	// CreateRelease
	r, _, err := c.CreateRelease(repo.Owner.UserName, repo.Name, CreateReleaseOption{
		TagName:      "awesome",
		Target:       "master",
		Title:        "Release 1",
		Note:         "yes it's awesome",
		IsDraft:      true,
		IsPrerelease: true,
	})
	assert.NoError(t, err)
	assert.EqualValues(t, "awesome", r.TagName)
	assert.EqualValues(t, true, r.IsPrerelease)
	assert.EqualValues(t, true, r.IsDraft)
	assert.EqualValues(t, "Release 1", r.Title)
	assert.EqualValues(t, fmt.Sprintf("%s/api/v1/repos/%s/releases/%d", c.url, repo.FullName, r.ID), r.URL)
	assert.EqualValues(t, "master", r.Target)
	assert.EqualValues(t, "yes it's awesome", r.Note)
	assert.EqualValues(t, c.username, r.Publisher.UserName)
	rl, _, _ = c.ListReleases(repo.Owner.UserName, repo.Name, ListReleasesOptions{})
	assert.Len(t, rl, 1)

	// GetRelease
	r2, _, err := c.GetRelease(repo.Owner.UserName, repo.Name, r.ID)
	assert.NoError(t, err)
	assert.EqualValues(t, r, r2)
	r2, _, err = c.GetReleaseByTag(repo.Owner.UserName, repo.Name, r.TagName)
	assert.NoError(t, err)
	assert.EqualValues(t, r, r2)
	// test fallback
	r2, _, err = c.fallbackGetReleaseByTag(repo.Owner.UserName, repo.Name, r.TagName)
	assert.NoError(t, err)
	assert.EqualValues(t, r, r2)

	// EditRelease
	r2, _, err = c.EditRelease(repo.Owner.UserName, repo.Name, r.ID, EditReleaseOption{
		Title:        "Release Awesome",
		Note:         "",
		IsDraft:      OptionalBool(false),
		IsPrerelease: OptionalBool(false),
	})
	assert.NoError(t, err)
	assert.EqualValues(t, r.Target, r2.Target)
	assert.EqualValues(t, false, r2.IsDraft)
	assert.EqualValues(t, false, r2.IsPrerelease)
	assert.EqualValues(t, r.Note, r2.Note)

	// DeleteRelease
	_, err = c.DeleteRelease(repo.Owner.UserName, repo.Name, r.ID)
	assert.NoError(t, err)
	rl, _, _ = c.ListReleases(repo.Owner.UserName, repo.Name, ListReleasesOptions{})
	assert.Len(t, rl, 0)

	// CreateRelease
	_, _, err = c.CreateRelease(repo.Owner.UserName, repo.Name, CreateReleaseOption{
		TagName: "aNewReleaseTag",
		Target:  "master",
		Title:   "Title of aNewReleaseTag",
	})
	assert.NoError(t, err)

	// DeleteReleaseByTag
	_, err = c.DeleteReleaseByTag(repo.Owner.UserName, repo.Name, "aNewReleaseTag")
	assert.NoError(t, err)
	rl, _, _ = c.ListReleases(repo.Owner.UserName, repo.Name, ListReleasesOptions{})
	assert.Len(t, rl, 0)
	_, err = c.DeleteReleaseByTag(repo.Owner.UserName, repo.Name, "aNewReleaseTag")
	assert.Error(t, err)

	// Test Response if try to get not existing release
	_, resp, err := c.GetRelease(repo.Owner.UserName, repo.Name, 1234)
	assert.Error(t, err)
	if assert.NotNil(t, resp) {
		assert.EqualValues(t, 404, resp.StatusCode)
	}
	_, resp, err = c.GetReleaseByTag(repo.Owner.UserName, repo.Name, "not_here")
	assert.Error(t, err)
	if assert.NotNil(t, resp) {
		assert.EqualValues(t, 404, resp.StatusCode)
	}
	_, resp, err = c.fallbackGetReleaseByTag(repo.Owner.UserName, repo.Name, "not_here")
	assert.Error(t, err)
	if assert.NotNil(t, resp) {
		assert.EqualValues(t, 404, resp.StatusCode)
	}
}
