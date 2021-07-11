// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package types

// AuthConfig contains authorization information for connecting to a cloud service.
type AuthConfig struct {
	Platform        Platform `json:"platform"`
	RegionID        string   `json:"region_id,omitempty"`
	SecretID        string   `json:"secret_id,omitempty"`
	SecretKey       string   `json:"secret_key,omitempty"`
	AccountID       string   `json:"account_id,omitempty"`
	AccessKeyID     string   `json:"access_key_id,omitempty"`
	AccessKeySecret string   `json:"access_key_secret,omitempty"`
}
