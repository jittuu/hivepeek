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
	Events      []*Event
	PreviousUrl string
	NextUrl     string
	League      string
	Season      string
	Date        time.Time
	IsAdmin     bool
}

func (e *Event) RatingDiff() int {
	if e.HRating == 0 || e.HNetRating == 0 || e.ARating == 0 || e.ANetRating == 0 {
		return 0
	}

	h := float64(e.HRating) + float64(e.HNetRating)
	a := float64(e.ARating) + float64(e.ANetRating)
	diff := h / (h + a) * 100
	if diff > 100 {
		diff = 100
	}

	return int(diff+0.5) - 50
}

func (e *Event) ResultString() string {
	switch {
	case e.HGoal > e.AGoal:
		return "W"
	case e.HGoal == e.AGoal:
		return "D"
	case e.HGoal < e.AGoal:
		return "L"
	}

	return ""
}

func (e *Event) ResultClass() string {
	switch e.ResultString() {
	case "W":
		return "win"
	case "D":
		return "draw"
	case "L":
		return "lose"
	}

	return ""
}

func (e *Event) RatingDiffClass() string {
	diff := e.RatingDiff()
	switch {
	case diff > 0:
		return "positive"
	case diff < 0:
		return "negative"
	}

	return ""
}

func (e *Event) Calculated() bool {
	if e.HRating != 0 || e.HNetRating != 0 || e.ARating != 0 || e.ANetRating != 0 {
		return true
	}

	return false
}
