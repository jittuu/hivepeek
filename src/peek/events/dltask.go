package events

import (
	"appengine"
	"encoding/csv"
	"io"
	"peek/ds"
)

type dlTask struct {
	context        appengine.Context
	season, league string
	w              io.Writer
}

func (t *dlTask) getAllEvents() ([]*Event, error) {
	dst, keys, err := ds.GetAllEvents(t.context, t.league, t.season)
	if err != nil {
		return nil, err
	}

	events := make([]*Event, len(dst))
	for i, e := range dst {
		events[i] = &Event{
			Event: e,
			Id:    keys[i].IntID(),
		}
	}

	startTime := func(e1, e2 *Event) bool {
		return e1.StartTime.Before(e2.StartTime)
	}

	By(startTime).Sort(events)

	return events, nil
}

func (t *dlTask) exec() error {
	events, err := t.getAllEvents()
	if err != nil {
		return err
	}

	records := CsvEvents(events).Csv()

	csvw := csv.NewWriter(t.w)
	return csvw.WriteAll(records)
}
