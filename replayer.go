package recmd

import (
	"sort"
	"time"
)

type ReplayEventType int

const (
	StdOut ReplayEventType = iota
	StdErr
	StdIn
)

type ReplayEvent struct {
	Data            []byte
	DataType        ReplayEventType
	Delay           time.Duration
	DelaySinceStart time.Duration
}

type Replayer struct {
	C       <-chan *ReplayEvent
	started bool
}

func NewReplayer(record Record) *Replayer {
	return newReplayer(record, false, false)
}

func NewQuickReplayer(record Record) *Replayer {
	return newReplayer(record, true, false)
}

func NewBackgroundReplayer(record Record) *Replayer {
	return newReplayer(record, false, true)
}

func newReplayer(record Record, quick bool, background bool) *Replayer {
	c := make(chan *ReplayEvent)

	events := []*ReplayEvent{}

	pipes := map[ReplayEventType]map[time.Duration][]byte{
		StdOut: record.StdOut(),
		StdErr: record.StdErr(),
		StdIn:  record.StdIn(),
	}

	for eventType, pipe := range pipes {
		for k, v := range pipe {
			events = append(events, &ReplayEvent{
				Data:            v,
				DataType:        eventType,
				Delay:           0,
				DelaySinceStart: k,
			})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].DelaySinceStart < events[j].DelaySinceStart
	})

	for i, event := range events {
		if i == 0 {
			event.Delay = event.DelaySinceStart
			continue
		}

		event.Delay = event.DelaySinceStart - events[i-1].DelaySinceStart
	}

	go func() {
		last := time.Now()
		for _, event := range events {

			if !quick && !background {
				<-time.After(event.Delay)
			}

			if background {
				<-time.After(event.Delay - time.Since(last))
				last = time.Now()
			}

			c <- event
		}

		close(c)
	}()

	return &Replayer{
		C:       c,
		started: false,
	}

}
