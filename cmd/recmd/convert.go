package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/scaxyz/recmd"
	"github.com/urfave/cli/v2"
)

func ConvertToStr(c *cli.Context) error {
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

	record := &recmd.PlainableRecord{
		Record: &recmd.Record{},
	}

	err = json.Unmarshal(jsonData, record)
	if err != nil {
		return err
	}

	jsonWithStrings, err := json.Marshal(&record)
	if err != nil {
		return err
	}

	outputPath := c.Args().Get(1)

	if strings.TrimSpace(outputPath) == "" {
		ext := filepath.Ext(recordFile)
		basenameAndPath := strings.TrimSuffix(recordFile, ext)
		newFilePath := fmt.Sprint(basenameAndPath, ".plain", ext)
		outputPath = newFilePath
	}

	outputFile, err := os.Create(outputPath)

	if err != nil {
		return err
	}

	outputFile.Write(jsonWithStrings)

	return nil

}
