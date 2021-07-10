package main

import (
	"os"

	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Raika/internal/cmd"
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
	}
	if err := app.Run(os.Args); err != nil {
		log.Error("%v", err)
	}
}
