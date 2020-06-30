// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGlobalSettings(t *testing.T) {
	log.Println("== TestGetGlobalSettings ==")
	c := newTestClient()

	uiSettings, err := c.GetGlobalUISettings()
	assert.NoError(t, err)
	expectedAllowedReactions := []string{"+1", "-1", "laugh", "hooray", "confused", "heart", "rocket", "eyes"}
	assert.ElementsMatch(t, expectedAllowedReactions, uiSettings.AllowedReactions)

	repoSettings, err := c.GetGlobalRepoSettings()
	assert.NoError(t, err)
	assert.False(t, repoSettings.HTTPGitDisabled)
	assert.False(t, repoSettings.MirrorsDisabled)
}
