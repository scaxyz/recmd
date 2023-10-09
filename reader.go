package recmd

import (
	"io"
	"sort"
	"time"
)

type RecordReader struct {
	data             map[time.Duration][]byte
	sortedTimePoints []time.Duration
	index            int
	readCount        int
	ignoreTime       bool
}

func NewReader(record Record) io.Reader {

	data := make(map[time.Duration][]byte)

	pipes := []map[time.Duration][]byte{
		record.StdOut(),
		record.StdIn(),
		record.StdErr(),
	}

	for _, pipe := range pipes {
		for k, v := range pipe {
			data[k] = v
		}
	}

	// Collect time durations
	timePoints := make([]time.Duration, len(data))
	i := 0
	for k := range data {
		timePoints[i] = k
		i++
	}

	// Sort time durations in ascending order
	sort.Slice(timePoints, func(i, j int) bool {
		return timePoints[i] < timePoints[j]
	})

	// Return new RecordReader object
	return &RecordReader{
		data:             data,
		sortedTimePoints: timePoints,
	}

}

func (rr *RecordReader) Reset() {
	rr.index = 0
	rr.readCount = 0
}

// IgnoreTime sets the ignoreTime field of the RecordReader struct to true.
//
// This function does not take any parameters and does not return any values.
func (rr *RecordReader) IgnoreTime() {
	rr.ignoreTime = true
}

func (rr *RecordReader) RespectTime() {
	rr.ignoreTime = false
}

// Read reads data from the RecordReader into the provided byte slice.
//
// It returns the number of bytes read and an error if any.
func (rr *RecordReader) Read(p []byte) (n int, err error) {
	// Check if there is no data left to read
	if len(rr.sortedTimePoints[rr.index:]) == 0 {
		return 0, io.EOF
	}

	// Get the current time point
	timePoint := rr.sortedTimePoints[rr.index]

	// Calculate the time difference between the current time point and the previous one
	var diff time.Duration = 0
	if rr.index != 0 {
		diff = timePoint - rr.sortedTimePoints[rr.index-1]
	}

	// Get the data to read from the current time point and readCount
	data := rr.data[timePoint][rr.readCount:]

	// Copy the data into the provided byte slice
	n = copy(p, data)

	// Update the readCount
	rr.readCount += n

	// Check if all data from the current time point has been read
	if rr.readCount == len(rr.data[timePoint]) {
		// Reset the readCount and move to the next time point
		rr.readCount = 0
		rr.index++
	}

	// Check if its the first read for the time point and it has a diff and the ignoreTime field is not set
	if (rr.readCount == 0 && diff != 0) && !rr.ignoreTime {
		<-time.NewTimer(diff).C
	}

	// Return the number of bytes read and nil error
	return n, nil
}
