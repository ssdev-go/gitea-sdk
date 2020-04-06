// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOauth2(t *testing.T) {
	log.Println("== TestOauth2Application ==")
	c := newTestClient()

	user := createTestUser(t, "oauth2_user", c)
	c.SetSudo(user.UserName)

	newApp, err := c.CreateOauth2(CreateOauth2Option{Name: "test", RedirectURIs: []string{"http://test/test",}})
	assert.NoError(t, err)
	assert.NotNil(t, newApp)
	assert.EqualValues(t, "test", newApp.Name)

	a, err := c.ListOauth2(ListOauth2Option{})
	assert.NoError(t, err)
	assert.Len(t, a, 1)
	assert.EqualValues(t, newApp.Name, a[0].Name)

	assert.NoError(t, c.DeleteOauth2(newApp.ID))
}
