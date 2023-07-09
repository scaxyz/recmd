package recmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

type JsonStrRecord struct {
	Record *Record `json:",inline"`
}

type RecordReader struct {
	data       map[time.Duration][]byte
	timepoints []time.Duration
	index      int
	readCount  int
	ignoreTime bool
}

// IgnoreTime sets the ignoreTime field of the RecordReader struct to true.
//
// This function does not take any parameters and does not return any values.
func (rr *RecordReader) IgnoreTime() {
	rr.ignoreTime = true
}

// Reader returns a RecordReader object.
//
// It creates a map of time durations to byte slices from the Out, In, and Err fields of the Record object.
// It then sorts the time durations in ascending order and returns a new RecordReader object with the created map and the sorted time durations.
func (r *Record) Reader() *RecordReader {
	data := make(map[time.Duration][]byte)

	// Copy values from r.Out, r.In, and r.Err to data
	for k, v := range r.Out {
		data[k] = v
	}
	for k, v := range r.In {
		data[k] = v
	}
	for k, v := range r.Err {
		data[k] = v
	}

	// Collect time durations
	timepoints := make([]time.Duration, len(data))
	i := 0
	for k := range data {
		timepoints[i] = k
		i++
	}

	// Sort time durations in ascending order
	sort.Slice(timepoints, func(i, j int) bool {
		return timepoints[i] < timepoints[j]
	})

	// Return new RecordReader object
	return &RecordReader{
		data:       data,
		timepoints: timepoints,
	}
}

// Read reads data from the RecordReader into the provided byte slice.
//
// It returns the number of bytes read and an error if any.
func (rr *RecordReader) Read(p []byte) (n int, err error) {
	// Check if there is no data left to read
	if len(rr.timepoints[rr.index:]) == 0 {
		return 0, io.EOF
	}

	// Get the current timepoint
	timepoint := rr.timepoints[rr.index]

	// Calculate the time difference between the current timepoint and the previous one
	var diff time.Duration = 0
	if rr.index != 0 {
		diff = timepoint - rr.timepoints[rr.index-1]
	}

	// Get the data to read from the current timepoint and readCount
	data := rr.data[timepoint][rr.readCount:]

	// Copy the data into the provided byte slice
	n = copy(p, data)

	// Update the readCount
	rr.readCount += n

	// Check if all data from the current timepoint has been read
	if rr.readCount == len(rr.data[timepoint]) {
		// Reset the readCount and move to the next timepoint
		rr.readCount = 0
		rr.index++
	}

	// Check if its the first read for the timepoint and it has a diff and the ignoreTime field is not set
	if (rr.readCount == 0 && diff != 0) && !rr.ignoreTime {
		<-time.NewTimer(diff).C
	}

	// Return the number of bytes read and nil error
	return n, nil
}

func (sjr *JsonStrRecord) MarshalJSON() ([]byte, error) {

	withBytes, err := json.Marshal(sjr.Record)
	if err != nil {
		return nil, err
	}
	asMap := make(map[string]interface{})
	err = json.Unmarshal(withBytes, &asMap)

	if err != nil {
		return nil, err
	}

	withStrs := recursivMap(asMap, "", func(data map[string]interface{}, k, path string) {
		if path == ".Command" {
			return
		}

		if maybeBase64, ok := data[k].(string); ok {
			decoded, err := base64.StdEncoding.DecodeString(maybeBase64)

			if err != nil {
				return
			}

			data[k] = string(decoded)
		}

	})
	return json.Marshal(withStrs)

}

func recursivMap(data map[string]interface{}, path string, callback func(data map[string]interface{}, k string, path string)) map[string]interface{} {
	if data == nil {
		return data
	}

	for k, v := range data {
		nextPath := fmt.Sprint(path, ".", k)

		if v, ok := v.(map[string]interface{}); ok {
			recursivMap(v, nextPath, callback)
			continue
		}
		if v, ok := v.([]interface{}); ok {
			for _, item := range v {
				if item, ok := item.(map[string]interface{}); ok {
					recursivMap(item, nextPath, callback)
				}
			}
			continue
		}

		callback(data, k, nextPath)

	}
	return data

}

func (sjr *JsonStrRecord) UnmarshalJSON(data []byte) error {
	asStringMap := make(map[string]interface{})
	err := json.Unmarshal(data, &asStringMap)
	if err != nil {
		return err
	}

	mapWithBytes := recursivMap(asStringMap, "", func(data map[string]interface{}, k, path string) {
		if path == ".Command" {
			return
		}

		if maybeBase64, ok := data[k].(string); ok {
			decoded, err := base64.StdEncoding.DecodeString(maybeBase64)
			if err == nil {
				data[k] = decoded
				return
			}
			data[k] = []byte(maybeBase64)
		}

	})

	jsonWithBytes, err := json.Marshal(mapWithBytes)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonWithBytes, sjr.Record)
	if err != nil {
		return err
	}

	return nil
}
