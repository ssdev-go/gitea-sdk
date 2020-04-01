// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import (
	"encoding/base64"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileCreateUpdateGet(t *testing.T) {
	log.Println("== TestFileCreateUpdateGet ==")
	c := newTestClient()

	repo, err := createTestRepo(t, "ChangeFiles", c)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	raw, err := c.GetFile(repo.Owner.UserName, repo.Name, "master", "README.md")
	assert.NoError(t, err)
	assert.EqualValues(t, "IyBDaGFuZ2VGaWxlcwoKQSB0ZXN0IFJlcG86IENoYW5nZUZpbGVz", base64.StdEncoding.EncodeToString(raw))

	newFile, err := c.CreateFile(repo.Owner.UserName, repo.Name, "A", CreateFileOptions{
		FileOptions: FileOptions{
			Message: "create file A",
		},
		Content: "ZmlsZUEK",
	})
	assert.NoError(t, err)
	raw, _ = c.GetFile(repo.Owner.UserName, repo.Name, "master", "A")
	assert.EqualValues(t, "ZmlsZUEK", base64.StdEncoding.EncodeToString(raw))

	updatedFile, err := c.UpdateFile(repo.Owner.UserName, repo.Name, "A", UpdateFileOptions{
		DeleteFileOptions: DeleteFileOptions{
			FileOptions: FileOptions{
				Message: "add a new line",
			},
			SHA: newFile.Content.SHA,
		},
		Content: "ZmlsZUEKCmFuZCBhIG5ldyBsaW5lCg==",
	})
	assert.NoError(t, err)
	assert.NotNil(t, updatedFile)

	file, err := c.GetContents(repo.Owner.UserName, repo.Name, "master", "A")
	assert.NoError(t, err)
	assert.EqualValues(t, updatedFile.Content.SHA, file.SHA)
	assert.EqualValues(t, &updatedFile.Content.Content, &file.Content)
}
