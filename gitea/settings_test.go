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

	uiSettings, _, err := c.GetGlobalUISettings()
	assert.NoError(t, err)
	expectedAllowedReactions := []string{"+1", "-1", "laugh", "hooray", "confused", "heart", "rocket", "eyes"}
	assert.ElementsMatch(t, expectedAllowedReactions, uiSettings.AllowedReactions)

	repoSettings, _, err := c.GetGlobalRepoSettings()
	assert.NoError(t, err)
	assert.EqualValues(t, &GlobalRepoSettings{
		HTTPGitDisabled: false,
		MirrorsDisabled: false,
	}, repoSettings)

	apiSettings, _, err := c.GetGlobalAPISettings()
	assert.NoError(t, err)
	assert.EqualValues(t, &GlobalAPISettings{
		MaxResponseItems:       50,
		DefaultPagingNum:       30,
		DefaultGitTreesPerPage: 1000,
		DefaultMaxBlobSize:     10485760,
	}, apiSettings)

	attachSettings, _, err := c.GetGlobalAttachmentSettings()
	assert.NoError(t, err)
	assert.EqualValues(t, &GlobalAttachmentSettings{
		Enabled:      true,
		AllowedTypes: ".docx,.gif,.gz,.jpeg,.jpg,.log,.pdf,.png,.pptx,.txt,.xlsx,.zip",
		MaxSize:      4,
		MaxFiles:     5,
	}, attachSettings)
}
