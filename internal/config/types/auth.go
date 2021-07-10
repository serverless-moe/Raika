// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package types

import (
	"github.com/wuhan005/Raika/internal/platform"
)

// AuthConfig contains authorization information for connecting to a cloud service.
type AuthConfig struct {
	Platform    platform.PlatformType `json:"platform"`
	AccessToken string                `json:"access_token,omitempty"`
	AccessKey   string                `json:"access_key,omitempty"`
}
