package events

import (
	"appengine"
	"fmt"
	"io"
	"math"
	"peek/ds"
	"time"
)

const betAmt = 100

type runTask struct {
	context  appengine.Context
	w        io.Writer
	season   string
	league   string
	diff     float64
	minPrice float64
	maxPrice float64
}

type MonthResults []*MonthResult

type MonthResult struct {
	Date   time.Time
	Bets   int
	Profit float64
	Events []*EventResult
}

type EventResult struct {
	*Event
	Profit float64
}

func (t *runTask) exec() (MonthResults, error) {
	events, err := t.getAllEvents()
	if err != nil {
		return nil, err
	}

	months := make(map[time.Time]*MonthResult)
	playedMatchCount := make(map[int64]int)

	for _, e := range events {
		fmt.Fprintf(t.w, "calculate for %s vs %s \n", e.Home, e.Away)

		if !e.Calculated() {
			break
		}

		fmt.Fprintf(t.w, "played count: (H: %d) (A: %d) \n", playedMatchCount[e.HomeId], playedMatchCount[e.AwayId])
		if playedMatchCount[e.HomeId] > 5 && playedMatchCount[e.AwayId] > 5 {
			y, m, _ := e.StartTime.Date()
			mdate := time.Date(y, m, 1, 0, 0, 0, 0, e.StartTime.Location())
			result := months[mdate]
			if result == nil {
				result = &MonthResult{Date: mdate}
				months[mdate] = result
			}

			if profit := t.evalEvent(e); profit != 0 {
				result.Bets += 1
				result.Profit += profit
				result.Events = append(result.Events, &EventResult{Event: e, Profit: profit})
			}
		}

		playedMatchCount[e.HomeId] += 1
		playedMatchCount[e.AwayId] += 1
	}
	fmt.Fprintf(t.w, "months count: %d \n", len(months))

	result := make(MonthResults, 0)
	for _, mr := range months {
		fmt.Fprintf(t.w, "bets: %d, profit: %.2f \n", mr.Bets, mr.Profit)
		result = append(result, mr)
	}

	return result, nil
}

func (t *runTask) evalEvent(e *Event) float64 {
	fmt.Fprintln(t.w, "start eval")
	diff := e.RatingDiff()
	fmt.Fprintf(t.w, "rating diff: %+.2f \n", diff)
	if diff > 0 {
		fmt.Fprintf(t.w, "betting for home to win with %.2f \n", e.MaxOdds.Home)
		// home to win
		if diff > t.diff && e.MaxOdds.Home >= t.minPrice && e.MaxOdds.Home <= t.maxPrice {
			// have bet
			if e.HGoal > e.AGoal {
				fmt.Fprintln(t.w, "WON")
				return betAmt * (e.MaxOdds.Home - 1)
			} else {
				fmt.Fprintln(t.w, "LOST")
				return betAmt * -1
			}
		}
	} else if diff < 0 {
		lay := e.LayPrice(e.MaxOdds.Home)
		fmt.Fprintf(t.w, "betting for home NOT to win with %.2f \n", lay)
		// home NOT to win
		if math.Abs(diff) > t.diff && lay >= t.minPrice && lay <= t.maxPrice {
			// have bet
			if e.HGoal > e.AGoal {
				fmt.Fprintln(t.w, "LOST")
				return betAmt * -1
			} else {
				fmt.Fprintln(t.w, "WON")
				return betAmt * (lay - 1)
			}
		}
	}

	return 0.0
}

func (t *runTask) getAllEvents() ([]*Event, error) {
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

func (months MonthResults) Profit() (bets int, profit float64) {
	for _, m := range months {
		bets += m.Bets
		profit += m.Profit
	}

	return
}

func (mr *MonthResult) Won() bool {
	return mr.Profit > 0
}

func (er *EventResult) Won() bool {
	return er.Profit > 0
}
