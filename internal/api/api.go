// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/pkg/errors"
)

func Stop() error {
	// TODO shutdown gently.
	_, _ = request(http.MethodPost, "/stop")
	return nil
}

func Reload() error {
	resp, err := request(http.MethodPost, "/reload")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusNoContent {
		return errors.Wrapf(err, "unexpected status code %d: %v", resp.StatusCode, resp.ToString())
	}
	return nil
}
