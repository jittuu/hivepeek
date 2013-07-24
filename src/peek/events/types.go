package events

import (
	"peek/ds"
	"time"
)

type Event struct {
	*ds.Event
	Id int64
}

type Team struct {
	*ds.Team
	Id int64
}

type GameWeek struct {
	Events      []*ds.Event
	PreviousUrl string
	NextUrl     string
	League      string
	Season      string
	Date        time.Time
}
