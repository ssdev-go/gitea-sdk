// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

import "fmt"

type searchUsersResponse struct {
	Users []*User `json:"data"`
}

// SearchUsers finds users by query
func (c *Client) SearchUsers(query string, limit int) ([]*User, error) {
	resp := new(searchUsersResponse)
	err := c.getParsedResponse("GET", fmt.Sprintf("/users/search?q=%s&limit=%d", query, limit), nil, nil, &resp)
	return resp.Users, err
}
