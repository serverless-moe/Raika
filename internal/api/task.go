// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package api

import (
	"net/http"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

func RunTask(functionName string) error {
	resp, err := request(http.MethodPost, "/task/run?functionName="+functionName)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status code %d: %v", resp.StatusCode, resp.ToString())
	}
	log.Trace("Response: %q", resp.ToString())
	return nil
}
