// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"time"

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
				&cli.StringFlag{Name: "description", Usage: "Function description", Required: false},
				&cli.Int64Flag{Name: "memory", Usage: "Function runtime memory size", Required: true},
				&cli.IntFlag{Name: "init-timeout", Usage: "Function runtime initialization timeout", Required: true},
				&cli.IntFlag{Name: "runtime-timeout", Usage: "Function runtime timeout", Required: true},
				&cli.StringFlag{Name: "binary-file", Usage: "Function binary file", Required: true},
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
			//client := tencentcloud.New(platform.AuthenticateOptions{
			//	tencentcloud.RegionIDField:  p.RegionID,
			//	tencentcloud.SecretIDField:  p.SecretID,
			//	tencentcloud.SecretKeyField: p.SecretKey,
			//})
			//platforms = append(platforms, client)
		default:
			return errors.Errorf("unsupported platform: %q", p)
		}
	}

	name := c.String("name")
	binaryFile := c.String("binary-file")
	description := c.String("description")
	memorySize := c.Int64("memory")
	initTimeout := c.Int("init-timeout")
	runtimeTimeout := c.Int("runtime-timeout")

	for _, p := range platforms {
		triggerURL, err := p.CreateFunction(platform.CreateFunctionOptions{
			Name:                  name,
			Description:           description,
			MemorySize:            memorySize,
			Environment:           map[string]string{},
			InitializationTimeout: time.Duration(initTimeout) * time.Second,
			RuntimeTimeout:        time.Duration(runtimeTimeout) * time.Second,
			HTTPPort:              9000, // For tencentcloud
			File:                  binaryFile,
		})
		if err != nil {
			log.Error("Failed to create function on %s: %v", p.Name(), err)
			continue
		}

		log.Info("[ %s ] - %s", p.Name(), triggerURL)
	}
	return nil
}
