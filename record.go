package recmd

import (
	"fmt"
	"io"
	"time"

	"github.com/samber/lo"
)

type RecordFormat string

const (
	FormatString RecordFormat = "string"
	FormatBase64 RecordFormat = "base64"
)

type Record interface {
	Format() RecordFormat
	Command() string
	StdOut() map[time.Duration][]byte
	StdIn() map[time.Duration][]byte
	StdErr() map[time.Duration][]byte
	Reader() io.Reader
	ConvertTo(format RecordFormat) (Record, error)
	ExitCode() int
}

type ByteRecord struct {
	JsonFormat RecordFormat `json:"format"`

	Cmd   string                   `json:"command"`
	Out   map[time.Duration][]byte `json:"out"`
	In    map[time.Duration][]byte `json:"in"`
	Err   map[time.Duration][]byte `json:"err"`
	ExitC int                      `json:"exitcode"`
}

type StringRecord struct {
	JsonFormat RecordFormat `json:"format"`

	Cmd   string                   `json:"command"`
	Out   map[time.Duration]string `json:"out"`
	In    map[time.Duration]string `json:"in"`
	Err   map[time.Duration]string `json:"err"`
	ExitC int                      `json:"exitcode"`
}

// Reader returns a RecordReader object.
//
// It creates a map of time durations to byte slices from the Out, In, and Err fields of the Record object.
// It then sorts the time durations in ascending order and returns a new RecordReader object with the created map and the sorted time durations.
func (r *ByteRecord) Reader() io.Reader {
	return NewReader(r)
}

func (br *ByteRecord) StdOut() map[time.Duration][]byte {
	return br.Out
}

func (br *ByteRecord) StdIn() map[time.Duration][]byte {
	return br.In
}

func (br *ByteRecord) StdErr() map[time.Duration][]byte {
	return br.Err
}

func (br *ByteRecord) Command() string {
	return br.Cmd
}

func (br *ByteRecord) ExitCode() int {
	return br.ExitC
}

func (br *ByteRecord) Format() RecordFormat {
	if br.JsonFormat == "" {
		br.JsonFormat = FormatBase64
	}
	return br.JsonFormat
}

func (br *ByteRecord) ConvertTo(format RecordFormat) (Record, error) {
	switch format {
	case FormatString:
		convert := func(byteMap map[time.Duration][]byte) map[time.Duration]string {
			convertedMap := make(map[time.Duration]string)
			for k, v := range byteMap {
				convertedMap[k] = string(v)
			}
			return convertedMap
		}
		return &StringRecord{
			Cmd:        br.Cmd,
			Out:        convert(br.Out),
			In:         convert(br.In),
			Err:        convert(br.Err),
			JsonFormat: FormatString,
		}, nil
	case FormatBase64:
		return &ByteRecord{
			Cmd:        br.Cmd,
			Out:        cloneMap(br.Out),
			In:         cloneMap(br.In),
			Err:        cloneMap(br.Err),
			JsonFormat: br.JsonFormat,
			ExitC:      br.ExitC,
		}, nil
	default:
		return nil, fmt.Errorf("unknown format: %s", format)
	}

}

func (sr *StringRecord) Reader() io.Reader {
	return NewReader(sr)
}

func (sr *StringRecord) StdOut() map[time.Duration][]byte {
	return lo.MapValues[time.Duration, string, []byte](sr.Out, func(value string, _ time.Duration) []byte {
		return []byte(value)
	})
}

func (sr *StringRecord) StdIn() map[time.Duration][]byte {
	return lo.MapValues[time.Duration, string, []byte](sr.In, func(value string, _ time.Duration) []byte {
		return []byte(value)
	})
}

func (sr *StringRecord) StdErr() map[time.Duration][]byte {
	return lo.MapValues[time.Duration, string, []byte](sr.Err, func(value string, _ time.Duration) []byte {
		return []byte(value)
	})
}

func (sr *StringRecord) Command() string {
	return sr.Cmd
}

func (sr *StringRecord) ExitCode() int {
	return sr.ExitC
}

func (sr *StringRecord) Format() RecordFormat {
	if sr.JsonFormat == "" {
		sr.JsonFormat = FormatString
	}
	return sr.JsonFormat
}

func (sr *StringRecord) ConvertTo(format RecordFormat) (Record, error) {
	switch format {
	case FormatString:
		return &StringRecord{
			Cmd:        sr.Cmd,
			Out:        cloneMap(sr.Out),
			In:         cloneMap(sr.In),
			Err:        cloneMap(sr.Err),
			JsonFormat: sr.JsonFormat,
			ExitC:      sr.ExitC,
		}, nil
	case FormatBase64:
		convert := func(stringMap map[time.Duration]string) map[time.Duration][]byte {
			convertedMap := make(map[time.Duration][]byte)
			for k, v := range stringMap {
				convertedMap[k] = []byte(v)
			}
			return convertedMap

		}
		return &ByteRecord{
			Cmd:        sr.Cmd,
			Out:        convert(sr.Out),
			In:         convert(sr.In),
			Err:        convert(sr.Err),
			JsonFormat: FormatBase64,
			ExitC:      sr.ExitC,
		}, nil
	default:
		return nil, fmt.Errorf("unknown format: %s", format)
	}
}

func cloneMap[K comparable, V any](originalMap map[K]V) map[K]V {
	clonedMap := make(map[K]V)
	for key, value := range originalMap {
		clonedMap[key] = value
	}
	return clonedMap
}
