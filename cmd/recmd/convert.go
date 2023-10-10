package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/scaxyz/recmd"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v2"
)

func ConvertToStr(ctx *cli.Context) error {
	recordFile := ctx.Args().First()

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

	err = json.Unmarshal(jsonData, &record)
	if err != nil {
		return err
	}
	strRecord, err := record.ConvertTo(recmd.FormatString)
	if err != nil {
		return err
	}

	outputPath := ctx.Args().Get(1)

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
	defer outputFile.Close()

	err = json.NewEncoder(outputFile).Encode(strRecord)
	if err != nil {
		return err
	}

	return nil

}
