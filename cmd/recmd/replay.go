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

	reader := record.Reader()

	if ctx.Bool("no-delays") {
		if r, ok := reader.(interface{ IgnoreTime() }); ok {
			r.IgnoreTime()
		}
	}

	buffer := make([]byte, replayBufferSize)

	for {
		data := buffer
		n, err := reader.Read(data)
		data = data[:n]
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fmt.Print(string(data))
	}
	exitCode := record.ExitCode()
	if ctx.IsSet("exit-code") {
		exitCode = ctx.Int("exit-code")
	}

	os.Exit(exitCode)

	return nil
}
