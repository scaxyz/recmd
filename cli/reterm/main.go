package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/scaxyz/recmd"
	"github.com/urfave/cli/v2"
)

func main() {

	tmpRecordFile := "record.json"

	app := cli.NewApp()

	app.Commands = []*cli.Command{
		{
			Name:  "rec",
			Usage: "Record something",
			Action: func(c *cli.Context) error {
				fmt.Println("Recording...")

				recorder := recmd.NewRecorder()

				record, err := recorder.Record("bash", os.Stdin)
				fmt.Printf("err: %#v\n", err)
				fmt.Printf("record: %#v\n", record)

				file, err := os.Create(tmpRecordFile)

				if err != nil {
					panic(err)
				}
				defer file.Close()

				recordJSON, err := json.Marshal(record)
				if err != nil {
					panic(err)
				}

				_, err = file.Write(recordJSON)
				if err != nil {
					panic(err)
				}

				return nil
			},
		},
		{
			Name:  "replay",
			Usage: "Replay something",
			Action: func(c *cli.Context) error {
				fmt.Println("Replaying...")

				file, err := os.Open(tmpRecordFile)
				if err != nil {
					panic(err)
				}
				defer file.Close()
				jsonData, err := io.ReadAll(file)

				if err != nil {
					panic(err)
				}

				record := &recmd.Record{}

				err = json.Unmarshal(jsonData, record)

				if err != nil {
					panic(err)
				}

				reader := record.Reader()
				buffer := make([]byte, 1024)
				fmt.Println(record.Command)
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
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}

}
