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

func Replay(c *cli.Context) error {

	recordFile := c.Args().First()

	file, err := os.Open(recordFile)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
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

	fmt.Println("Replaying: ", record.Command())

	reader := record.Reader()

	if c.Bool("no-delays") {
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

	return nil
}
