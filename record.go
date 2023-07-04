package recmd

import (
	"io"
	"sort"
	"time"
)

type Record struct {
	Command string
	Out     map[time.Duration][]byte
	In      map[time.Duration][]byte
	Err     map[time.Duration][]byte
}

type RecordReader struct {
	data        map[time.Duration][]byte
	timespoints []time.Duration
	index       int
	readCount   int
	ignoreTime  bool
}

func (rr *RecordReader) IgnoreTime() {
	rr.ignoreTime = true
}

func (r *Record) Reader() *RecordReader {
	data := make(map[time.Duration][]byte)
	for k, v := range r.Out {
		data[k] = v
	}
	for k, v := range r.In {
		data[k] = v
	}
	for k, v := range r.Err {
		data[k] = v
	}

	timepoints := make([]time.Duration, len(data))
	i := 0
	for k := range data {
		timepoints[i] = k
		i++
	}

	sort.Slice(timepoints, func(i, j int) bool {
		return timepoints[i] < timepoints[j]
	})

	return &RecordReader{
		data:        data,
		timespoints: timepoints,
	}
}

// Read reads data from the RecordReader into the provided byte slice.
//
// It returns the number of bytes read and an error if any.
func (rr *RecordReader) Read(p []byte) (n int, err error) {

	// return if nothing left
	if len(rr.timespoints[rr.index:]) == 0 {
		return 0, io.EOF
	}

	timepoint := rr.timespoints[rr.index]
	var diff time.Duration = 0
	if rr.index != 0 {
		diff = timepoint - rr.timespoints[rr.index-1]
	}

	// reset time to dont wait mutilble times if the buffer is smaller than the data
	wait := rr.readCount == 0 && diff != 0

	data := rr.data[timepoint][rr.readCount:]

	n = copy(p, data)
	rr.readCount += n

	if rr.readCount == len(rr.data[timepoint]) {
		rr.readCount = 0
		rr.index++
	}

	// procced to the next data if the current has been read

	if wait && !rr.ignoreTime {
		<-time.NewTimer(diff).C
	}

	return n, nil

}
