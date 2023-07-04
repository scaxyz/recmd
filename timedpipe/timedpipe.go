package timedpipe

import (
	"fmt"
	"io"
	"time"
)

type Pipe struct {
	dataWrite map[time.Duration][]byte
	dataRead  map[time.Duration][]byte
	start     time.Time
	output    io.Writer
	input     io.Reader
}

type PipeOption func(*Pipe)

func WithOutput(w io.Writer) PipeOption {
	return func(t *Pipe) {
		t.SetOutput(w)
	}
}

func WithInput(r io.Reader) PipeOption {
	return func(t *Pipe) {
		t.SetInput(r)
	}
}

func StartNow() PipeOption {
	return func(t *Pipe) {
		t.SetStartTime(time.Now())
	}
}

func New(options ...PipeOption) *Pipe {
	pipe := &Pipe{
		dataWrite: make(map[time.Duration][]byte),
		dataRead:  make(map[time.Duration][]byte),
	}
	for _, option := range options {
		option(pipe)
	}
	return pipe
}

func (t *Pipe) Write(p []byte) (n int, err error) {

	if t.start.IsZero() {
		t.start = time.Now()
	}

	copied := make([]byte, len(p))
	copy(copied, p)

	t.dataWrite[time.Since(t.start)] = copied

	if t.output != nil {
		return t.output.Write(p)
	}

	return len(p), nil
}

func (t *Pipe) Read(p []byte) (n int, err error) {

	if t.start.IsZero() {
		t.start = time.Now()
	}

	if t.input != nil {

		n, err = t.input.Read(p)
		copied := make([]byte, n)
		copy(copied, p)
		t.dataRead[time.Since(t.start)] = copied
		return
	}

	return 0, fmt.Errorf("no input")

}

func (t *Pipe) GetReadData() map[time.Duration][]byte {
	return t.dataRead
}

func (t *Pipe) GetWriteData() map[time.Duration][]byte {
	return t.dataWrite
}

func (t *Pipe) SetInput(r io.Reader) *Pipe {
	t.input = r
	return t
}

func (t *Pipe) SetOutput(w io.Writer) *Pipe {
	t.output = w
	return t
}

func (t *Pipe) SetStartTime(time time.Time) *Pipe {
	t.start = time
	return t
}
