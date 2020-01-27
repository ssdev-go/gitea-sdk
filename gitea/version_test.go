// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	log.Printf("== TestVersion ==")
	c := newTestClient()
	rawVersion, err := c.ServerVersion()
	assert.NoError(t, err)
	assert.True(t, true, rawVersion != "")

	assert.NoError(t, c.CheckServerVersionConstraint(">= 1.11.0"))
	assert.Error(t, c.CheckServerVersionConstraint("< 1.11.0"))
}
