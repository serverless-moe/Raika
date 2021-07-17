// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package types

import (
	"time"
)

// Function represent as an individual function.
type Function struct {
	PlatformID string    `json:"platform_id"`
	URL        string    `json:"string"`
	CreatedAt  time.Time `json:"created_at"`

	Name                  string            `json:"name"`
	Description           string            `json:"description"`
	MemorySize            int64             `json:"memory_size"`
	Environment           map[string]string `json:"environment"`
	InitializationTimeout time.Duration     `json:"initialization_timeout"`
	RuntimeTimeout        time.Duration     `json:"runtime_timeout"`
	HTTPPort              int               `json:"http_port"`
	File                  string            `json:"file"`
}
