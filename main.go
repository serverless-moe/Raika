package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/cmd"
	"github.com/wuhan005/Raika/internal/config"
	"github.com/wuhan005/Raika/internal/store"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "Raika"
	app.Commands = []*cli.Command{
		cmd.Daemon,
		cmd.Platform,
		cmd.Function,
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "config-file", Value: config.DefaultConfigPath, Usage: "Config file path"},
		&cli.StringFlag{Name: "function-file", Value: store.DefaultFunctionPath, Usage: "Function file path"},
		&cli.StringFlag{Name: "task-file", Value: store.DefaultTaskPath, Usage: "Task file path"},
	}
	app.Before = func(c *cli.Context) error {
		if err := store.Functions.Init(c.String("function-file")); err != nil {
			return errors.Wrap(err, "load function file")
		}
		if err := store.Tasks.Init(c.String("task-file")); err != nil {
			return errors.Wrap(err, "load task file")
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Error("%v", err)
	}
}
