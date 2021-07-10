// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/config"
	"github.com/wuhan005/Raika/internal/platform"
	"github.com/wuhan005/Raika/internal/platform/aliyun"
	"github.com/wuhan005/Raika/internal/types"
)

var Function = &cli.Command{
	Name:  "function",
	Usage: "Manage the functions",
	Subcommands: []*cli.Command{
		{
			Name:   "create",
			Usage:  "Create a new function to the cloud service",
			Action: createFunction,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "name", Usage: "Function name", Required: true},
				&cli.StringFlag{Name: "template", Usage: "Function template file", Required: true},
			},
		},
	},
}

func createFunction(c *cli.Context) error {
	configFilePath := c.String("config-file")
	configFile := config.New(configFilePath)
	if err := configFile.Load(); err != nil {
		return errors.Wrap(err, "load config file")
	}

	platforms := make([]platform.Cloud, 0, len(configFile.AuthConfigs))
	for _, p := range configFile.AuthConfigs {
		switch p.Platform {
		case types.Aliyun:
			client := aliyun.New(platform.AuthenticateOptions{
				aliyun.RegionIDField:        p.RegionID,
				aliyun.AccountIDField:       p.AccountID,
				aliyun.AccessKeyIDField:     p.AccessKeyID,
				aliyun.AccessKeySecretField: p.AccessKeySecret,
			})
			platforms = append(platforms, client)
		case types.TencentCloud:

		default:
			return errors.Errorf("unsupported platform: %q", p)
		}
	}

	for _, p := range platforms {
		err := p.CreateFunction(platform.CreateFunctionOptions{})
		if err != nil {
			log.Error("Failed to create function on %s: %v", p.Name(), err)
		}
	}
	return nil
}
