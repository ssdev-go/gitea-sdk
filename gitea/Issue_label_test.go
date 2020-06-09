// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLabels test label related func
func TestLabels(t *testing.T) {
	log.Println("== TestLabels ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "LabelTestsRepo", c)
	assert.NoError(t, err)

	createOpts := CreateLabelOption{
		Name:        " ",
		Description: "",
		Color:       "",
	}
	err = createOpts.Validate()
	assert.Error(t, err)
	assert.EqualValues(t, "invalid color format", err.Error())
	createOpts.Color = "12345f"
	err = createOpts.Validate()
	assert.Error(t, err)
	assert.EqualValues(t, "empty name not allowed", err.Error())
	createOpts.Name = "label one"

	label1, err := c.CreateLabel(repo.Owner.UserName, repo.Name, createOpts)
	assert.NoError(t, err)
	assert.EqualValues(t, createOpts.Name, label1.Name)
	assert.EqualValues(t, createOpts.Color, label1.Color)

}
