// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package types

type Platform string

const (
	Aliyun       Platform = "aliyun"
	TencentCloud Platform = "tencentcloud"
	AWS          Platform = "aws"
)

func (p Platform) Check() bool {
	switch p {
	case Aliyun, TencentCloud, AWS:
		return true
	}
	return false
}
