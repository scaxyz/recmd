package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var replayBufferSize = 1024 * 1024
var defaultFileTimeFormat = "20060102_150405"
var defaultOutputTemplate = "recmd-{{ .CmdBaseName }}-{{ .Time }}.json"
var outputTemplateOnTemplateError = fmt.Sprintf("recmd-%s-template-error.json", time.Now().Format(defaultFileTimeFormat))

var version = "development"

func main() {

	app := cli.NewApp()
	app.Version = version
	app.Usage = "record or replay inputs and outputs of a command"

	app.Commands = []*cli.Command{
		{
			Name:    "record",
			Aliases: []string{"rec"},
			Usage:   "Records the following command",
			Flags: []cli.Flag{
				&cli.PathFlag{
					Name:    "input",
					Usage:   "Use file as stdin",
					Aliases: []string{"i", "in", "if"},
				},
				&cli.PathFlag{
					Name:    "output",
					Usage:   "Output file",
					Aliases: []string{"o", "out", "of"},
					Value:   defaultOutputTemplate,
				},
				&cli.BoolFlag{
					Name:    "save-with-plain-text",
					Aliases: []string{"plain-text", "plain", "pt", "p"},
					Usage:   "Saves to json with 'in','out' and 'err' as plain texts instead of base64 encodings",
				},
				&cli.StringFlag{
					Name:  "time-format",
					Usage: "time format for the output template, accessible with {{ .Time }}",
					Value: defaultFileTimeFormat,
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
				&cli.IntFlag{
					Name:    "exit-code",
					Usage:   "Overwrites the exit-code from the replay",
					Aliases: []string{"code", "ec"},
				},
				&cli.BoolFlag{
					Name:    "no-delays",
					Usage:   "Ignore delays while replaying",
					Aliases: []string{"quick"},
				},
			},
			Action: Replay,
		},
		{
			Name:      "convert-to-plain-text",
			Aliases:   []string{"conv-plain", "cpt"},
			Usage:     "Converts an record with 'in', 'out' and 'error' as base64 to one which uses plain text instead, (default-output: <input-name>-string.<input-ext>)",
			Action:    ConvertToStr,
			UsageText: "recmd convert-to-plain-text <input-file> [output-file]",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
