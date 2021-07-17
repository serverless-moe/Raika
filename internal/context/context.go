// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"encoding/json"
	"net/http"

	"github.com/flamego/flamego"
	log "unknwon.dev/clog/v2"
)

type Context struct {
	flamego.Context
}

func (c *Context) JSON(statusCode int, data interface{}) {
	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(statusCode)
	if err := json.NewEncoder(c.ResponseWriter()).Encode(data); err != nil {
		log.Error("Failed to encode response body: %v", err)
	}
}

func (c *Context) NoContent() {
	c.ResponseWriter().WriteHeader(http.StatusNoContent)
}

func (c *Context) Success(data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"error": false,
		"msg":   "success",
		"data":  data,
	})
}

func (c *Context) Error(statusCode int, message string) {
	c.JSON(statusCode, map[string]interface{}{
		"error": true,
		"msg":   message,
	})
}

func (c *Context) ServerError() {
	c.Error(http.StatusInternalServerError, "internal server error")
}

func Contexter() flamego.Handler {
	return func(ctx flamego.Context) {
		c := Context{
			Context: ctx,
		}

		c.Map(c)
	}
}
