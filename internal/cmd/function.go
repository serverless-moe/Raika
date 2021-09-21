// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/config"
	"github.com/wuhan005/Raika/internal/platform"
	"github.com/wuhan005/Raika/internal/platform/aliyun"
	"github.com/wuhan005/Raika/internal/platform/aws"
	"github.com/wuhan005/Raika/internal/platform/tencentcloud"
	"github.com/wuhan005/Raika/internal/store"
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
				&cli.StringSliceFlag{Name: "platform", Usage: "Platform to deploy", Required: false},
				&cli.StringSliceFlag{Name: "env", Usage: "Environment variables", Required: false},
				&cli.StringFlag{Name: "trigger", Usage: "Function trigger method", Required: false, DefaultText: "http"},
				&cli.StringFlag{Name: "cron", Usage: "Cron expression for timer trigger", Required: false, DefaultText: "0 30 * * * *"},
			},
		},
		{
			Name:   "list",
			Usage:  "List all the functions",
			Action: listFunctions,
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
	platformNames := c.StringSlice("platform")
	platformNameSet := make(map[types.Platform]struct{})
	for _, platformName := range platformNames {
		platformNameSet[types.Platform(platformName)] = struct{}{}
	}

	for _, p := range configFile.AuthConfigs {
		_, ok := platformNameSet[p.Platform]
		if len(platformNameSet) != 0 && !ok {
			continue
		}

		switch p.Platform {
		case types.Aliyun:
			client := aliyun.New(platform.AuthenticateOptions{
				"id":                        fmt.Sprintf("%s@%s@%s", types.Aliyun, p.AccountID, p.RegionID),
				aliyun.RegionIDField:        p.RegionID,
				aliyun.AccountIDField:       p.AccountID,
				aliyun.AccessKeyIDField:     p.AccessKeyID,
				aliyun.AccessKeySecretField: p.AccessKeySecret,
			})
			platforms = append(platforms, client)
		case types.TencentCloud:
			client := tencentcloud.New(platform.AuthenticateOptions{
				"id":                        fmt.Sprintf("%s@%s@%s", types.TencentCloud, p.SecretID, p.RegionID),
				tencentcloud.RegionIDField:  p.RegionID,
				tencentcloud.SecretIDField:  p.SecretID,
				tencentcloud.SecretKeyField: p.SecretKey,
			})
			platforms = append(platforms, client)
		case types.AWS:
			client := aws.New(platform.AuthenticateOptions{
				"id":               fmt.Sprintf("%s@%s@%s", types.AWS, p.AccountID, p.RegionID),
				aws.RegionIDField:  p.RegionID,
				aws.AccountIDField: p.AccountID,
				aws.AccessKeyField: p.AccessKeyID,
				aws.SecretKeyField: p.SecretKey,
			})
			platforms = append(platforms, client)
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
	environmentVariables := c.StringSlice("env")
	trigger := c.String("trigger")
	cron := c.String("cron")

	envs := make(map[string]string)
	// Parse environment variables.
	for _, env := range environmentVariables {
		kv := strings.SplitN(env, "=", 2)
		if len(kv) != 2 {
			continue
		}
		envs[kv[0]] = kv[1]
	}

	for _, p := range platforms {
		log.Info("Create function %q on %s", name, p)

		opts := platform.CreateFunctionOptions{
			Name:                  name,
			Description:           description,
			MemorySize:            memorySize,
			EnvironmentVariables:  envs,
			InitializationTimeout: time.Duration(initTimeout) * time.Second,
			RuntimeTimeout:        time.Duration(runtimeTimeout) * time.Second,
			File:                  binaryFile,

			TriggerType: trigger,
			CronString:  cron,
			HTTPPort:    9000, // For tencentcloud
		}
		triggerURL, err := p.CreateFunction(opts)
		if err != nil {
			log.Error("Failed to create function on %s: %v", p, err)
			continue
		}
		// Save the function into file.
		if err := store.Functions.Set(name, p.GetID(), triggerURL, opts); err != nil {
			log.Error("Failed to save function to file: %v", err)
		}

		log.Info("[ %s ] - %s", p, triggerURL)
	}
	return nil
}

func listFunctions(c *cli.Context) error {
	for name, platforms := range store.Functions.Functions {
		log.Info("-  %s", name)
		for _, p := range platforms {
			log.Trace("   [%s] %s", p.PlatformID, p.URL)
		}
	}
	return nil
}
