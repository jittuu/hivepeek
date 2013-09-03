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

func (e *Event) RatingDiff() float64 {
	if !e.Calculated() {
		return 0.0
	}

	h := e.HRating + ((e.HNetRating + e.HFormRating) / 2)
	a := e.ARating + ((e.ANetRating + e.AFormRating) / 2)
	hPer := h / (h + a) * 100
	aPer := 100 - hPer
	diff := hPer - aPer
	return diff
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

func (e *Event) HasEnoughRating() bool {
	return e.HFormRatingLen == 5 && e.AFormRatingLen == 5
}

func (e *Event) LayPrice(back float64) float64 {
	return (1 / (back - 1)) + 1
}

type Events []*Event

func (events Events) Find(home, away string) *Event {
	for _, e := range events {
		if e.Home == home && e.Away == away {
			return e
		}
	}

	return nil
}

type RunTaskResult struct {
	Profit  float64
	Bets    int
	Results MonthResults
}

func (r *RunTaskResult) Won() bool {
	return r.Profit > 0
}
