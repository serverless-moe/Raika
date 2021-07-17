// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package daemon

import (
	"math/rand"
	"net/http"

	"github.com/pkg/errors"

	"github.com/wuhan005/Raika/internal/store"
)

func runFunction(functionName string) (*http.Response, error) {
	for fn := range store.Tasks.Tasks {
		if fn == functionName {
			platformFunctions, err := store.Functions.Get(functionName)
			if err != nil {
				return nil, errors.Wrapf(err, "get platform function: %q", functionName)
			}

			index := rand.Intn(len(platformFunctions))
			platformFunction := platformFunctions[index]
			return http.Get(platformFunction.URL)
		}
	}
	return nil, errors.Wrapf(store.ErrFunctionNotExists, "platform function: %q", functionName)
}
