package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/scaxyz/recmd"
	"github.com/urfave/cli/v2"
)

func Record(ctx *cli.Context) error {
	commands := ctx.Args().Slice()
	now := time.Now()
	if len(commands) == 0 {
		return fmt.Errorf("no command specified")
	}

	var input io.Reader = nil

	if ctx.Bool("interactive") {
		input = os.Stdin
	}

	if ctx.Path("input") != "" {
		fileInput, err := os.Open(ctx.Path("input"))
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

	var finalRecord recmd.Record = record
	if ctx.Bool("save-with-plain-text") {
		finalRecord, err = finalRecord.ConvertTo(recmd.FormatString)
		if err != nil {
			return err
		}
	}

	recordJSON, err := json.Marshal(finalRecord)
	if err != nil {
		return err
	}

	outputFilePath := buildOutputFilePath(finalRecord, ctx.Path("output"), now.Format(ctx.String("time-format")))

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.Write(recordJSON)
	if err != nil {
		return err
	}

	log.Println("wrote recording to " + outputFilePath)

	return nil
}

func buildOutputFilePath(record recmd.Record, templateStr string, time string) string {

	outputTemplate, err := template.New("output").Parse(templateStr)
	if err != nil {
		log.Println("error: parsing output path template: " + err.Error())
		return outputTemplateOnTemplateError
	}

	builder := strings.Builder{}

	err = outputTemplate.Execute(&builder, NewTemplateContext(record, time))
	if err != nil {
		log.Println("error: executing output path template: " + err.Error())
		return outputTemplateOnTemplateError
	}

	return builder.String()

}
