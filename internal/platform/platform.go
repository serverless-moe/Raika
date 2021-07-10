// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package platform

type AuthenticateOptions map[string]string

type Cloud interface {
	Authenticate() error
}
