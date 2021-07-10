// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package aliyun

import (
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/platform"
)

func (c *Client) CreateFunction(opts platform.CreateFunctionOptions) error {
	service, err := c.GetRaikaService()
	if err == ErrRaikaServiceNotFound {
		log.Trace("Raika service not found on aliyun, create...")
		service, err = c.CreateService("Raika-service", "Service for Raika.")
		if err != nil {
			return errors.Wrap(err, "create service")
		}
	} else if err != nil {
		return err
	}

	log.Trace("%+v", service)

	return nil
}
