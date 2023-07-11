package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var replayBufferSize = 1024 * 1024
var defaultFileTimeFormat = "20060102_150405"

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
					Aliases: []string{"i", "in", "if"},
				},
				&cli.PathFlag{
					Name:    "output",
					Usage:   "Output file",
					Aliases: []string{"o", "out", "of"},
					Value:   fmt.Sprintf("recmd-%s.json", time.Now().Format(defaultFileTimeFormat)),
				},
				&cli.BoolFlag{
					Name:    "save-with-plain-text",
					Aliases: []string{"plain-text", "plain", "pt", "p"},
					Usage:   "Saves to json with 'in','out' and 'err' as plain texts instead of base64 encodings",
				},
				&cli.BoolFlag{
					Name:    "interactive",
					Aliases: []string{"inter", "stdin"},
					Usage:   "Use standard input",
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
		{
			Name:    "convert-to-plain-text",
			Aliases: []string{"conv-plain", "cpt"},
			Usage:   "Converts an record with 'in', 'out' and 'error' as base64 to on which uses plain text instead",
			Action:  ConvertToStr,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
