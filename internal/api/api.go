// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"
)

func Stop() error {
	// TODO shutdown gently.
	_, _ = request(http.MethodPost, "/stop")
	return nil
}

func Reload() error {
	_, _ = request(http.MethodPost, "/reload")
	return nil
}
