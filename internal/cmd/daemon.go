// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

var Daemon = &cli.Command{
	Name:   "daemon",
	Usage:  "Run Raika daemon",
	Action: runDaemon,
	Flags: []cli.Flag{

	},
}

func runDaemon(c *cli.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	return nil
}
