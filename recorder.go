package recmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/scaxyz/recmd/timedpipe"
)

type Recorder struct{}

// NewRecorder creates a new Recorder.
//
// It takes no parameters and returns a pointer to a Recorder.
func NewRecorder() *Recorder {
	return &Recorder{}
}

// Record records a command and returns a Record object with the command's output and error.
//
// cmdStr: the command to be executed.
// Returns a pointer to the Record object and an error.
func (r *Recorder) Record(cmdStr string, input io.Reader, args ...string) (*Record, error) {

	if cmdStr == "" {
		return nil, fmt.Errorf("empty command")
	}

	cmd := exec.Command(cmdStr, args...)

	return r.RecordCmd(cmd, input)
}

func (*Recorder) RecordCmd(cmd *exec.Cmd, input io.Reader) (*Record, error) {

	if cmd == nil {
		return nil, fmt.Errorf("empty command")
	}

	errP := timedpipe.New(timedpipe.WithOutput(os.Stderr), timedpipe.StartNow())
	outP := timedpipe.New(timedpipe.WithOutput(os.Stdout), timedpipe.StartNow())
	inP := timedpipe.New(timedpipe.WithInput(input), timedpipe.StartNow())

	cmd.Stderr = errP
	cmd.Stdout = outP

	if input != nil {
		cmd.Stdin = inP
	}

	record := &Record{
		Command: cmd.String(),
		Out:     outP.GetWriteData(),
		In:      inP.GetReadData(),
		Err:     errP.GetWriteData(),
	}
	err := cmd.Run()
	return record, err
}
