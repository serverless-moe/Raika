// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package types

import (
	"time"
)

type Task struct {
	FunctionName string        `json:"function_name"`
	Duration     time.Duration `json:"duration"`
	Enabled      bool          `json:"enabled"`
}
