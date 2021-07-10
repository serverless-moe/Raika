package main

import (
	"os"

	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/cmd"
	"github.com/wuhan005/Raika/internal/config"
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
	}

	if err := app.Run(os.Args); err != nil {
		log.Error("%v", err)
	}
}
