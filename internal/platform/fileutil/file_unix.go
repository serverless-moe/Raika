// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// +build !windows

package fileutil

import (
	"os"
	"syscall"
)

// CopyFilePermissions copies file ownership and permissions from "src" to "dst",
// ignoring any error during the process.
func CopyFilePermissions(src, dst string) {
	var (
		mode     os.FileMode = 0600
		uid, gid int
	)

	fi, err := os.Stat(src)
	if err != nil {
		return
	}
	if fi.Mode().IsRegular() {
		mode = fi.Mode()
	}
	if err := os.Chmod(dst, mode); err != nil {
		return
	}

	uid = int(fi.Sys().(*syscall.Stat_t).Uid)
	gid = int(fi.Sys().(*syscall.Stat_t).Gid)

	if uid > 0 && gid > 0 {
		_ = os.Chown(dst, uid, gid)
	}
}
