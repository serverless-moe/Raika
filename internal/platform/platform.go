// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package platform

import (
	"github.com/wuhan005/Raika/internal/types"
)

type AuthenticateOptions map[string]string

type Cloud interface {
	String() string
	Platform() types.Platform
	GetID() string
	Authenticate() error
	CreateFunction(opts CreateFunctionOptions) (string, error)
}
