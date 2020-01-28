// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyUser(t *testing.T) {
	log.Println("== TestMyUser ==")
	c := newTestClient()
	user, err := c.GetMyUserInfo()
	assert.NoError(t, err)

	assert.EqualValues(t, 1, user.ID)
	assert.EqualValues(t, "test01", user.UserName)
	assert.EqualValues(t, "test01@gitea.io", user.Email)
	assert.EqualValues(t, "", user.FullName)
	assert.EqualValues(t, getGiteaURL()+"/user/avatar/test01/-1", user.AvatarURL)
	assert.Equal(t, true, user.IsAdmin)
}
