// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import "github.com/urfave/cli/v2"

func stringFlag(name, value, usage string, required ...bool) *cli.StringFlag {
	var _required bool
	if len(required) == 1 {
		_required = required[0]
	}

	return &cli.StringFlag{
		Name:     name,
		Value:    value,
		Usage:    usage,
		Required: _required,
	}
}

func intFlag(name string, value int, usage string, required ...bool) *cli.IntFlag {
	var _required bool
	if len(required) == 1 {
		_required = required[0]
	}
	
	return &cli.IntFlag{
		Name:     name,
		Value:    value,
		Usage:    usage,
		Required: _required,
	}
}
