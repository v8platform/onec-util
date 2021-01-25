package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
	"v8platform/onec-util/cmd"
)

// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {

	app := &cli.App{
		Name:     "onec-util",
		Version:  version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name: "Aleksey Khorev",
			},
		},
		Usage:     "Command line utilities for server 1S.Enterprise",
		UsageText: "onec-util command [command options] [arguments...]",
		Copyright: "(c) 2021 Khorevaa",
		//Description: "Command line utilities for server 1S.Enterprise",
	}

	for _, command := range cmd.Commands {
		app.Commands = append(app.Commands, command.Cmd())
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
