package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/scaxyz/recmd"
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

	jSRecord := recmd.PlainableRecord{
		Record: &recmd.Record{},
	}
	record := jSRecord.Record

	err = json.Unmarshal(jsonData, &jSRecord)
	if err != nil {
		return err
	}

	fmt.Println("Replaying: ", record.Command)

	reader := record.Reader()

	if c.Bool("no-delays") {
		reader.IgnoreTime()
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
