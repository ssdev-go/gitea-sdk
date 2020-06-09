// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gitea

// GlobalSettings represent the global settings of a gitea instance witch is exposed by API
type GlobalSettings struct {
	AllowedReactions []string
}

// GetGlobalSettings get all global settings witch are exposed by API
func (c *Client) GetGlobalSettings() (settings GlobalSettings, err error) {
	settings.AllowedReactions, err = c.GetSettingAllowedReactions()
	return
}

// GetSettingAllowedReactions return reactions witch are allowed on a instance
func (c *Client) GetSettingAllowedReactions() ([]string, error) {
	if err := c.CheckServerVersionConstraint(">=1.13.0"); err != nil {
		return nil, err
	}
	var reactions []string
	return reactions, c.getParsedResponse("GET", "/settings/allowed_reactions", jsonHeader, nil, &reactions)
}
