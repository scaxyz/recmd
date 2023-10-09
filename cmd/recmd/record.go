package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/scaxyz/recmd"
	"github.com/urfave/cli/v2"
)

func Record(c *cli.Context) error {
	commands := c.Args().Slice()

	if len(commands) == 0 {
		return fmt.Errorf("no command specified")
	}

	var input io.Reader = nil

	if c.Bool("interactive") {
		input = os.Stdin
	}

	if c.Path("input") != "" {
		fileInput, err := os.Open(c.Path("input"))
		if err != nil {
			return err
		}
		defer fileInput.Close()
		input = fileInput
	}

	fmt.Printf("Recording: '%s'\n", strings.Join(commands, " "))

	recorder := recmd.NewRecorder()

	record, err := recorder.Record(commands[0], input, commands[1:]...)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(c.String("output"))
	if err != nil {
		return err
	}
	defer outputFile.Close()

	var finalRecord recmd.Record = record
	if c.Bool("save-with-plain-text") {
		finalRecord, err = finalRecord.ConvertTo(recmd.FormatString)
		if err != nil {
			return err
		}
	}

	recordJSON, err := json.Marshal(finalRecord)
	if err != nil {
		return err
	}

	_, err = outputFile.Write(recordJSON)
	if err != nil {
		return err
	}

	return nil
}
