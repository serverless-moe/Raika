// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/wuhan005/Raika/internal/platform"
	"github.com/wuhan005/Raika/internal/types"
)

var _ platform.Cloud = (*Client)(nil)

type Client struct {
	id                   string
	regionID             string
	accountID, roleName  string
	accessKey, secretKey string
}

func New(opts platform.AuthenticateOptions) *Client {
	return &Client{
		id:        opts["id"],
		regionID:  opts[RegionIDField],
		roleName:  opts[RoleName],
		accountID: opts[AccountIDField],
		accessKey: opts[AccessKeyField],
		secretKey: opts[SecretKeyField],
	}
}

func (c *Client) String() string {
	return string(c.Platform())
}

func (c *Client) Platform() types.Platform {
	return types.AWS
}

func (c *Client) GetID() string {
	return c.id
}

func (c *Client) Authenticate() error {
	_, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(c.accessKey, c.secretKey, ""),
		Region:      &c.regionID,
	})
	if err != nil {
		return err
	}

	return nil
}
