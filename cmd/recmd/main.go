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
					Name:        "output",
					Usage:       "Output file",
					Aliases:     []string{"o", "out", "of"},
					DefaultText: fmt.Sprintf("recmd-%s.json", time.Now().Format(defaultFileTimeFormat)),
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
			Name:    "convert-to-string",
			Aliases: []string{"conv-str", "cs"},
			Usage:   "Converts an record with bytes/base64 to on which uses strings instead",
			Action:  ConvertToStr,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
