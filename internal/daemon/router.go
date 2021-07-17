// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package daemon

import (
	gocontext "context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/flamego/flamego"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/context"
	"github.com/wuhan005/Raika/internal/store"
)

func Run() error {
	c := cron.New()
	taskEntrySets := make(map[string]cron.EntryID)

	// Cron task.
	for _, task := range store.Tasks.Tasks {
		if !task.Enabled {
			continue
		}

		entryID, err := c.AddFunc(fmt.Sprintf("@every %ds", task.Duration/time.Second), func() {
			// TODO handle error
			_, _ = runFunction(task.FunctionName)
		})
		if err != nil {
			continue
		}

		taskEntrySets[task.FunctionName] = entryID
	}

	c.Start()

	f := flamego.Classic()
	f.Use(context.Contexter())
	server := http.Server{
		Addr:    "127.0.0.1:3000",
		Handler: f,
	}

	f.Group("/task", func() {
		f.Post("/run", func(ctx context.Context) {
			functionName := ctx.Request().URL.Query().Get("functionName")
			resp, err := runFunction(functionName)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
			ctx.JSON(http.StatusOK, string(body))
		})

		f.Post("/enable", func(ctx context.Context) {
			functionName := ctx.Request().URL.Query().Get("functionName")
			task, err := store.Tasks.Get(functionName)
			if err != nil {
				ctx.Error(http.StatusNotFound, err.Error())
				return
			}

			entryID, err := c.AddFunc(fmt.Sprintf("@every %ds", task.Duration/time.Second), func() {
				// TODO handle error
				_, _ = runFunction(task.FunctionName)
			})
			if err != nil {
				ctx.Error(http.StatusInternalServerError, err.Error())
			}
			taskEntrySets[task.FunctionName] = entryID
		})

		f.Post("/disable", func(ctx context.Context) {
			functionName := ctx.Request().URL.Query().Get("functionName")
			entryID, ok := taskEntrySets[functionName]
			if ok {
				c.Remove(entryID)
				delete(taskEntrySets, functionName)
				err := store.Tasks.Disable(functionName)
				if err != nil {
					ctx.Error(http.StatusInternalServerError, errors.Wrap(err, "disable tasks").Error())
					return
				}
			}
			ctx.NoContent()
		})
	})

	f.Post("/stop", func(ctx context.Context) {
		err := server.Shutdown(gocontext.Background())
		if err != nil {
			ctx.Error(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.NoContent()
	})

	f.Post("/reload", func(ctx context.Context) {
		if err := store.Functions.Load(); err != nil {
			ctx.Error(http.StatusInternalServerError, errors.Wrap(err, "reload functions file").Error())
			return
		}
		if err := store.Tasks.Load(); err != nil {
			ctx.Error(http.StatusInternalServerError, errors.Wrap(err, "reload tasks file").Error())
			return
		}
		ctx.NoContent()
	})

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Trace("Server closed.")
		return nil
	}
	return err
}
