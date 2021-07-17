// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"

	"github.com/wuhan005/Raika/internal/api"
	"github.com/wuhan005/Raika/internal/daemon"
)

var Daemon = &cli.Command{
	Name:  "daemon",
	Usage: "Configure the Raika daemon",
	Subcommands: []*cli.Command{
		{
			Name:   "start",
			Usage:  "Start the Raika daemon",
			Action: startDaemon,
		},
		{
			Name:   "run",
			Usage:  "Run the daemon on frontend",
			Action: runDaemon,
		},
		{
			Name:   "stop",
			Usage:  "Stop the Raika daemon",
			Action: stopDaemon,
		},
		{
			Name:   "reload",
			Usage:  "Reload the config",
			Action: reloadConfig,
		},
		{
			Name:        "cron",
			Usage:       "Set the cron task",
			Subcommands: cronCommands,
		},
	},
	Flags: []cli.Flag{

	},
}

func startDaemon(_ *cli.Context) error {
	cmd := exec.Command(os.Args[0], "daemon", "run")
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func runDaemon(_ *cli.Context) error {
	return daemon.Run()
}

func stopDaemon(_ *cli.Context) error {
	return api.Stop()
}

func reloadConfig(_ *cli.Context) error {
	return api.Reload()
}
