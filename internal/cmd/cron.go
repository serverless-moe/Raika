// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/api"
	"github.com/wuhan005/Raika/internal/store"
)

var cronCommands = []*cli.Command{
	{
		Name:   "list",
		Usage:  "List all the cron tasks",
		Action: listCornTask,
	},
	{
		Name:   "create",
		Usage:  "Create a new cron task",
		Action: createCornTask,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Usage: "Function name", Required: true},
			&cli.IntFlag{Name: "duration", Usage: "Duration time", Required: true},
		},
	},
	{
		Name:   "delete",
		Usage:  "Delete the cron task",
		Action: deleteCornTask,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Usage: "Function name", Required: true},
		},
	},
	{
		Name:   "run",
		Usage:  "Run the cron task immediately",
		Action: runCornTask,
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Usage: "Function name", Required: true},
		},
	},
}

func listCornTask(_ *cli.Context) error {
	for functionName, task := range store.Tasks.Tasks {
		var status = "ENABLED"
		if !task.Enabled {
			status = "DISABLED"
		}
		log.Trace("[%s] %s/s (%s)", functionName, task.Duration/time.Second, status)
	}
	return nil
}

func createCornTask(c *cli.Context) error {
	functionName := c.String("name")
	secondDuration := c.Int("duration")

	if _, err := store.Functions.Get(functionName); err != nil {
		return errors.Wrap(err, "get function")
	}

	err := store.Tasks.Upsert(store.CreateTaskOptions{
		FunctionName: functionName,
		Duration:     time.Duration(secondDuration) * time.Second,
	})
	if err != nil {
		return errors.Wrap(err, "create task")
	}

	if err := api.Reload(); err != nil {
		return errors.Wrapf(err, "reload")
	}
	return nil
}

func deleteCornTask(c *cli.Context) error {
	functionName := c.String("name")
	return store.Tasks.Delete(functionName)
}

func runCornTask(c *cli.Context) error {
	functionName := c.String("name")
	return api.RunTask(functionName)
}
