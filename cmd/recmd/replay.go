package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/scaxyz/recmd"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v2"
)

func Replay(ctx *cli.Context) error {

	recordFile := ctx.Args().First()

	file, err := os.Open(recordFile)
	if err != nil {
		return err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// no defer since we are using os.Exit at the and
	err = file.Close()
	if err != nil {
		return err
	}

	format := gjson.Get(string(jsonData), "format")
	var record recmd.Record

	switch format.String() {
	case string(recmd.FormatBase64):
		record = &recmd.ByteRecord{}
	case string(recmd.FormatString):
		record = &recmd.StringRecord{}
	default:
		return fmt.Errorf("unknown format: %s", format.String())
	}

	err = json.Unmarshal(jsonData, record)
	if err != nil {
		return err
	}

	if !ctx.Bool("pure") {
		fmt.Println("Replaying: ", record.Command())
	}

	var replayer *recmd.Replayer

	if ctx.Bool("no-delays") {
		replayer = recmd.NewQuickReplayer(record)
	} else {
		replayer = recmd.NewReplayer(record)
	}

	for event := range replayer.C {
		switch event.DataType {
		case recmd.StdOut:
			os.Stdout.Write(event.Data)
		case recmd.StdErr:
			os.Stderr.Write(event.Data)
		case recmd.StdIn:
			os.Stdout.Write(event.Data)
		}
	}

	exitCode := record.ExitCode()
	if ctx.IsSet("exit-code") {
		exitCode = ctx.Int("exit-code")
	}

	os.Exit(exitCode)

	return nil
}
