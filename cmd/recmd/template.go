package main

import (
	"path/filepath"
	"strings"

	"github.com/scaxyz/recmd"
)

type TemplateContext struct {
	Time        string
	CmdBaseName string
	Record      struct {
		Format   string
		Command  string
		ExitCode int
	}
}

func NewTemplateContext(record recmd.Record, time string) *TemplateContext {
	cmd, _, _ := strings.Cut(record.Command(), " ")

	return &TemplateContext{
		Time:        time,
		CmdBaseName: filepath.Base(cmd),
		Record: struct {
			Format   string
			Command  string
			ExitCode int
		}{
			Format:   string(record.Format()),
			Command:  record.Command(),
			ExitCode: record.ExitCode(),
		},
	}
}
