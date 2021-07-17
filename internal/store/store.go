// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package store

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var HomePath, _ = os.UserHomeDir()
var DefaultFunctionPath = filepath.Join(HomePath, "./.raika/functions.json")
var DefaultTaskPath = filepath.Join(HomePath, "./.raika/tasks.json")

var ErrFunctionNotExists = errors.New("function not found")
