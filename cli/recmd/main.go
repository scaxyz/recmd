package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var replayBufferSize = 1024 * 1024

func main() {

	app := cli.NewApp()

	app.Commands = []*cli.Command{
		{
			Name:    "record",
			Aliases: []string{"rec"},
			Usage:   "Records the following command",
			Flags: []cli.Flag{
				&cli.PathFlag{
					Name:    "input",
					Usage:   "Use file as input",
					Aliases: []string{"i", "in"},
				},
				&cli.PathFlag{
					Name:     "output",
					Usage:    "Output file",
					Aliases:  []string{"o", "out"},
					Required: true,
				},
				&cli.BoolFlag{
					Name:  "no-stdin",
					Usage: "Disable standard input",
				},
			},
			Action: Record,
		},
		{
			Name:    "replay",
			Aliases: []string{"rep"},
			Usage:   "Replay a recorded command",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "no-delays",
					Usage: "Ignore delays while replaying",
				},
			},
			Action: Replay,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
