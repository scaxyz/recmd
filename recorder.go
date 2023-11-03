package recmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/scaxyz/recmd/timed"
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
func (r *Recorder) Record(cmdStr string, input io.Reader, args ...string) (Record, error) {

	if cmdStr == "" {
		return nil, fmt.Errorf("empty command")
	}

	cmd := exec.Command(cmdStr, args...)
	record, err := r.RecordCmd(cmd, input)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func (*Recorder) RecordCmd(cmd *exec.Cmd, input io.Reader) (Record, error) {

	if cmd == nil {
		return nil, fmt.Errorf("empty command")
	}
	startTime := time.Now()

	errP := timed.New(timed.WithOutput(os.Stderr), timed.SetStartTime(startTime))
	outP := timed.New(timed.WithOutput(os.Stdout), timed.SetStartTime(startTime))
	inP := timed.New(timed.WithInput(input), timed.SetStartTime(startTime))

	cmd.Stderr = errP
	cmd.Stdout = outP

	if input != nil {
		cmd.Stdin = inP
	}

	record := &ByteRecord{
		Cmd:        cmd.String(),
		Out:        outP.GetWriteData(),
		In:         inP.GetReadData(),
		Err:        errP.GetWriteData(),
		JsonFormat: FormatBase64,
	}

	stopChan := make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	doneChan := make(chan error)

	go func() {
		err := cmd.Run()
		doneChan <- err
	}()

	var err error

	select {
	case interrupt := <-stopChan:
		cmd.Process.Signal(interrupt)
	case err = <-doneChan:
		break
	}

	record.ExitC = cmd.ProcessState.ExitCode()

	if _, ok := err.(*exec.ExitError); ok {
		err = nil
	}

	if err != nil {
		return nil, err
	}

	return record, nil
}
