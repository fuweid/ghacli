package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"github.com/fuweid/ghacli/commands/workflowruns"
	"github.com/fuweid/ghacli/commands/workflows"
)

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "ghacli"
	app.Description = "A Github Action commandline"
	app.Usage = "A Github Action commandline"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "owner",
			Usage:    "The account owner of the repository. The name is not case sensitive.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "repo",
			Usage:    "The name of the repository. The name is not case sensitive.",
			Required: true,
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "The Github token of caller",
			EnvVar: "GITHUB_TOKEN",
		},
	}
	app.Commands = append(app.Commands,
		workflows.Command,
		workflowruns.Command,
	)
	return app
}

func main() {
	app := newCliApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", app.Name, err)
		os.Exit(1)
	}
}
