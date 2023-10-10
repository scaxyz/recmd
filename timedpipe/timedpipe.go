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
	started   bool
	output    io.Writer
	input     io.Reader
}

type PipeOption func(*Pipe)

// WithOutput sets the output writer for the Pipe.
//
// w: the output writer to set.
// Returns: a PipeOption function.
func WithOutput(w io.Writer) PipeOption {
	return func(t *Pipe) {
		t.SetOutput(w)
	}
}

// WithInput returns a PipeOption function that sets the input of a Pipe.
//
// r is the input reader.
// PipeOption is a function that modifies a Pipe.
func WithInput(r io.Reader) PipeOption {
	return func(t *Pipe) {
		t.SetInput(r)
	}
}

// StartNow returns a PipeOption function that sets the start time of the Pipe to the current time.
//
// It takes no parameters and returns a PipeOption.
func StartNow() PipeOption {
	return func(t *Pipe) {
		t.SetStartTime(time.Now())
	}
}

// SetStartTime sets the start time of the Pipe.
//
// time: the start time to set.
// Returns: a PipeOption function.
func SetStartTime(time time.Time) PipeOption {
	return func(t *Pipe) {
		t.start = time
	}
}

// New creates a new instance of Pipe with the provided options.
//
// The options parameter is variadic and allows for configuration of the Pipe.
// It accepts zero or more PipeOption values.
//
// The function returns a pointer to the created Pipe.
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

func (t *Pipe) setStartTime() {
	t.start = time.Now()
	t.started = true
}

// Write writes the given byte slice to the Pipe.
// Storing data and the time passed since the start time of the Pipe.
//
// It takes a parameter p which is a byte slice.
// It returns the number of bytes written and an error, if any.
func (t *Pipe) Write(p []byte) (n int, err error) {

	if !t.started {
		t.setStartTime()
	}

	if p == nil {
		return 0, fmt.Errorf("no buffer")
	}

	stamp := time.Now()

	copied := make([]byte, len(p))
	copy(copied, p)

	go func() {
		t.dataWrite[stamp.Sub(t.start)] = copied
	}()

	if t.output != nil {
		return t.output.Write(p)
	}

	return len(p), nil
}

// Read reads data from the input into the given byte slice.
// Storing the time passed since the start time of the pipe
//
// It returns the number of bytes read and any error encountered.
func (t *Pipe) Read(p []byte) (n int, err error) {

	if !t.started {
		t.setStartTime()
	}

	if t.input == nil {
		return 0, fmt.Errorf("no buffer")
	}

	stamp := time.Now()

	n, err = t.input.Read(p)
	if err != nil {
		return n, err
	}

	copied := make([]byte, n)
	copy(copied, p)

	go func() {
		t.dataRead[stamp.Sub(t.start)] = copied
	}()

	return n, nil

}

// GetReadData returns the map of bytes read from the pipe, indexed by time duration.
//
// It does not take any parameters.
// It returns a map of time duration to a slice of bytes.
func (t *Pipe) GetReadData() map[time.Duration][]byte {
	return t.dataRead
}

// GetWriteData returns the data to be written by the Pipe.
//
// It returns a map with time durations as keys and byte slices as values.
func (t *Pipe) GetWriteData() map[time.Duration][]byte {
	return t.dataWrite
}

// SetInput assigns an input reader to the Pipe.
//
// r - the input reader to be assigned.
// Returns the modified Pipe.
func (t *Pipe) SetInput(r io.Reader) *Pipe {
	t.input = r
	return t
}

// SetOutput sets the output for the Pipe.
//
// It takes a parameter `w` of type `io.Writer` and returns a pointer to `Pipe`.
func (t *Pipe) SetOutput(w io.Writer) *Pipe {
	t.output = w
	return t
}

// SetStartTime sets the start time of the Pipe.
//
// time - the time to set as the start time.
// Returns a pointer to the Pipe.
func (t *Pipe) SetStartTime(time time.Time) *Pipe {
	t.start = time
	return t
}
